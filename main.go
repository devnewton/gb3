package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"time"

	"./tribune"
)

var listenAddress string
var verboseMode bool

func init() {
	flag.StringVar(&listenAddress, "listen", ":6666", "TCP address to listen on")
	flag.BoolVar(&verboseMode, "verbose", false, "Verbose logging")
}

type gb3 struct {
	coincoins map[*tribune.Coincoin]struct{}
	join      chan *tribune.Coincoin
	leave     chan *tribune.Coincoin
	forward   chan tribune.Post
	tribunes  map[string]*tribune.Tribune
	poll      chan *tribune.Tribune
	indexer   tribune.Indexer
	store     tribune.Store
}

func newGb3() *gb3 {
	return &gb3{
		coincoins: make(map[*tribune.Coincoin]struct{}),
		join:      make(chan *tribune.Coincoin),
		leave:     make(chan *tribune.Coincoin),
		forward:   make(chan tribune.Post),
		tribunes: map[string]*tribune.Tribune{
			"euromussels": &tribune.Tribune{Name: "euromussels", BackendURL: "https://faab.euromussels.eu/data/backend.tsv", PostURL: "https://faab.euromussels.eu/add.php", PostField: "message"},
			"sveetch":     &tribune.Tribune{Name: "sveetch", BackendURL: "http://sveetch.net/tribune/remote/tsv/", PostURL: "http://sveetch.net/tribune/post/tsv/?last_id=1", PostField: "content"},
			"moules":      &tribune.Tribune{Name: "moules", BackendURL: "http://moules.org/board/backend/tsv", PostURL: "http://moules.org/board/add.php?backend=tsv", PostField: "message"},
			"ototu":       &tribune.Tribune{Name: "ototu", BackendURL: "https://ototu.euromussels.eu/goboard/backend/tsv", PostURL: "https://ototu.euromussels.eu/goboard/post", PostField: "message"},
			"dlfp":        &tribune.Tribune{Name: "dlfp", BackendURL: "https://linuxfr.org/board/index.tsv", PostURL: "https://linuxfr.org/api/v1/board", PostField: "message", AuthentificationType: tribune.OAuth2Authentification},
		},
		poll:    make(chan *tribune.Tribune),
		indexer: tribune.NewIndexer(),
		store:   tribune.NewStore(),
	}
}

func (a *gb3) forwardLoop() {
	for {
		select {
		case c := <-a.join:
			a.coincoins[c] = struct{}{}
		case c := <-a.leave:
			delete(a.coincoins, c)
			close(c.Send)
		case post := <-a.forward:
			js, err := json.Marshal(post)
			if nil != err {
				log.Println(err)
				continue
			}
			for c := range a.coincoins {
				c.Send <- tribune.CoincoinMessage{Post: post, Data: js}
			}
		}
	}
}

func (a *gb3) sendStoredPostsTo(c *tribune.Coincoin) {
	for _, t := range a.tribunes {
		posts, err := a.store.RetrieveLatests(t.Name)
		if nil != err {
			log.Println(err)
			continue
		}
		for _, p := range posts {
			js, err := json.Marshal(p)
			if nil != err {
				log.Println(err)
				continue
			}
			c.Send <- tribune.CoincoinMessage{Post: p, Data: js}
		}
	}
}

func (a *gb3) handlePoll(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	c := tribune.NewCoincoin(w, flusher)
	a.join <- c
	defer func() { a.leave <- c }()
	go a.sendStoredPostsTo(c)
	c.WriteLoop()
}

func (a *gb3) handlePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if nil != err {
		log.Println(err)
		return
	}
	t := a.tribunes[r.PostFormValue("tribune")]
	if nil != t {
		t.Post(r)
		a.poll <- t
	}
}

func (a *gb3) handleSearch(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("query")
	results, err := a.indexer.Search(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	err = encoder.Encode(results)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (a *gb3) pollTribunes() {
	for _, t := range a.tribunes {
		err := a.pollTribune(t)
		if nil != err {
			log.Println(err)
			continue
		}
	}
}

func (a *gb3) pollTribune(t *tribune.Tribune) error {
	if verboseMode {
		log.Printf("Poll %s\n", t.Name)
	}
	posts, err := t.Poll()
	if nil != err {
		return err
	}
	for _, p := range posts {
		a.forward <- p
	}
	go a.store.Save(t.Name, posts)
	go a.indexer.Index(posts)
	return err
}

func (a *gb3) pollLoop() {
	tick := time.Tick(30 * time.Second)
	for {
		select {
		case t := <-a.poll:
			a.pollTribune(t)
		case <-tick:
			a.pollTribunes()
		}
	}
}

func main() {
	flag.Parse()
	a := newGb3()
	go a.forwardLoop()
	go a.pollLoop()
	http.HandleFunc("/api/poll", func(w http.ResponseWriter, r *http.Request) {
		a.handlePoll(w, r)
	})
	http.HandleFunc("/api/post", func(w http.ResponseWriter, r *http.Request) {
		a.handlePost(w, r)
	})
	http.HandleFunc("/api/search", func(w http.ResponseWriter, r *http.Request) {
		a.handleSearch(w, r)
	})
	http.Handle("/", http.FileServer(http.Dir("gc2")))
	log.Printf("Listen to %s\n", listenAddress)
	err := http.ListenAndServe(listenAddress, nil)
	if nil != err {
		log.Fatal(err)
	}
}
