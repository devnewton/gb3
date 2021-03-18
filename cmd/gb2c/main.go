package main

import (
	"embed"
	"encoding/json"
	"flag"
	"io/fs"
	"log"
	"net/http"
	"time"
)

//go:embed gc2
var gc2Content embed.FS

var listenAddress string
var verboseMode bool

func init() {
	flag.StringVar(&listenAddress, "listen", ":16666", "TCP address to listen on")
	flag.BoolVar(&verboseMode, "verbose", false, "Verbose logging")
}

type gb3 struct {
	join     chan *Coincoin
	leave    chan *Coincoin
	forward  chan forwardMessage
	tribunes map[string]*Tribune
	poll     chan *Tribune
	indexer  Indexer
	store    Store
}

type forwardMessage struct {
	Post
	*Coincoin
}

func newGb3() *gb3 {
	return &gb3{
		join:    make(chan *Coincoin),
		leave:   make(chan *Coincoin),
		forward: make(chan forwardMessage),
		tribunes: map[string]*Tribune{
			"euromussels": {Name: "euromussels", BackendURL: "https://faab.euromussels.eu/data/backend.tsv", PostURL: "https://faab.euromussels.eu/add.php", PostField: "message"},
			"sveetch":     {Name: "sveetch", BackendURL: "http://sveetch.net/tribune/remote/tsv/", PostURL: "http://sveetch.net/tribune/post/tsv/?last_id=1", PostField: "content"},
			"moules":      {Name: "moules", BackendURL: "http://moules.org/board/backend/tsv", PostURL: "http://moules.org/board/add.php?backend=tsv", PostField: "message"},
			"ototu":       {Name: "ototu", BackendURL: "https://ototu.euromussels.eu/goboard/backend/tsv", PostURL: "https://ototu.euromussels.eu/goboard/post", PostField: "message"},
			"dlfp":        {Name: "dlfp", BackendURL: "https://linuxfr.org/board/index.tsv", PostURL: "https://linuxfr.org/api/v1/board", PostField: "message", AuthentificationType: OAuth2Authentification},
		},
		poll:    make(chan *Tribune),
		indexer: NewIndexer(),
		store:   NewStore(),
	}
}

func (g *gb3) forwardLoop() {
	coincoins := make(map[*Coincoin]struct{})
	for {
		select {
		case c := <-g.join:
			coincoins[c] = struct{}{}
		case c := <-g.leave:
			delete(coincoins, c)
		case msg := <-g.forward:
			js, err := json.Marshal(msg.Post)
			if nil != err {
				log.Println(err)
				continue
			}
			if nil == msg.Coincoin {
				aliveCoincoins := make(map[*Coincoin]struct{})
				for c := range coincoins {
					if err := c.Write(msg.Post, js); nil == err {
						aliveCoincoins[c] = struct{}{}
					} else {
						if verboseMode {
							log.Printf("Error writing to coincoin : %s\n", err)
						}
					}
				}
				coincoins = aliveCoincoins
			} else {
				if err := msg.Coincoin.Write(msg.Post, js); nil != err {
					delete(coincoins, msg.Coincoin)
					log.Printf("Error writing to coincoin : %s\n", err)
				}
			}
		}
	}
}

func (g *gb3) sendStoredPostsTo(c *Coincoin) {
	for _, t := range g.tribunes {
		posts, err := g.store.RetrieveLatests(t.Name)
		if nil != err {
			log.Println(err)
			continue
		}
		for _, p := range posts {
			g.forward <- forwardMessage{Post: p, Coincoin: c}
		}
	}
}

func (g *gb3) handlePoll(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	c := NewCoincoin(w, flusher)
	g.join <- c
	defer func() { g.leave <- c }()
	g.sendStoredPostsTo(c)
	select {
	case <-r.Context().Done():
	case <-time.After(60 * time.Minute):
	}
}

func (g *gb3) handlePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if nil != err {
		log.Println(err)
		return
	}
	t := g.tribunes[r.PostFormValue("tribune")]
	if nil != t {
		t.Post(r)
		g.poll <- t
	}
}

func (g *gb3) handleSearch(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("query")
	results, err := g.indexer.Search(query)
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

func (g *gb3) pollTribunes() {
	for _, t := range g.tribunes {
		err := g.pollTribune(t)
		if nil != err {
			log.Println(err)
			continue
		}
	}
}

func (g *gb3) pollTribune(t *Tribune) error {
	if verboseMode {
		log.Printf("Poll %s\n", t.Name)
	}
	posts, err := t.Poll()
	if nil != err {
		return err
	}
	for _, p := range posts {
		g.forward <- forwardMessage{Post: p}
	}
	go g.store.Save(t.Name, posts)
	go g.indexer.Index(posts)
	return err
}

func (g *gb3) pollLoop() {
	tick := time.Tick(30 * time.Second)
	for {
		select {
		case t := <-g.poll:
			g.pollTribune(t)
		case <-tick:
			g.pollTribunes()
		}
	}
}

func main() {
	flag.Parse()
	g := newGb3()
	go g.forwardLoop()
	go g.pollLoop()
	RegisterLinuxfrAPI()
	http.HandleFunc("/api/poll", func(w http.ResponseWriter, r *http.Request) {
		g.handlePoll(w, r)
	})
	http.HandleFunc("/api/post", func(w http.ResponseWriter, r *http.Request) {
		g.handlePost(w, r)
	})
	http.HandleFunc("/api/search", func(w http.ResponseWriter, r *http.Request) {
		g.handleSearch(w, r)
	})
	http.HandleFunc("/api/totoz/search", TotozSearch)
	http.HandleFunc("/api/emoji/search", EmojiSearch)
	gc2fs, err := fs.Sub(gc2Content, "gc2")
	if nil != err {
		log.Fatal(err)
	}
	http.Handle("/", http.FileServer(http.FS(gc2fs)))
	log.Printf("Listen to %s\n", listenAddress)
	err = http.ListenAndServe(listenAddress, nil)
	if nil != err {
		log.Fatal(err)
	}
}
