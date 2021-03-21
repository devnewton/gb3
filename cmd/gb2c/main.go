package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

var verboseMode bool

func init() {
	switch os.Getenv("GB2C_VERBOSE") {
	case "1", "true", "TRUE":
		verboseMode = true
	default:
		verboseMode = false
	}
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
	g := &gb3{
		join:    make(chan *Coincoin),
		leave:   make(chan *Coincoin),
		forward: make(chan forwardMessage),
		tribunes: map[string]*Tribune{
			"euromussels": {Name: "euromussels", BackendURL: "https://faab.euromussels.eu/data/backend.tsv", PostURL: "https://faab.euromussels.eu/add.php", PostField: "message"},
			"sveetch":     {Name: "sveetch", BackendURL: "http://sveetch.net/tribune/remote/tsv/", PostURL: "http://sveetch.net/tribune/post/tsv/?last_id=1", PostField: "content"},
			"moules":      {Name: "moules", BackendURL: "http://moules.org/board/backend/tsv", PostURL: "http://moules.org/board/add.php?backend=tsv", PostField: "message"},
			"ototu":       {Name: "ototu", BackendURL: "https://ototu.euromussels.eu/goboard/backend/tsv", PostURL: "https://ototu.euromussels.eu/goboard/post", PostField: "message"},
			"dlfp":        {Name: "dlfp", BackendURL: "https://linuxfr.org/board/index.tsv", PostURL: "https://linuxfr.org/gb2c/v1/board", PostField: "message", AuthentificationType: OAuth2Authentification},
		},
		poll:    make(chan *Tribune),
		indexer: NewIndexer(),
		store:   NewStore(),
	}
	gb0Tribunes := os.Getenv("GB2C_GB0_TRIBUNES")
	if len(gb0Tribunes) == 0 {
		gb0Tribunes = "gb0local:http://localhost:16667"
	}
	for _, gb0Tribune := range strings.Split(gb0Tribunes, ",") {
		gb0TribuneSplitted := strings.SplitN(gb0Tribune, ":", 2)
		if len(gb0TribuneSplitted) != 2 {
			fmt.Printf("Invalid gb0 tribune description in GB2C_GB0_TRIBUNES: %s\n", gb0Tribune)
			continue
		}
		gb0TribuneName := gb0TribuneSplitted[0]
		gb0TribuneURL := gb0TribuneSplitted[1]
		g.tribunes[gb0TribuneName] = &Tribune{Name: gb0TribuneName, BackendURL: gb0TribuneURL + "/gb0/tsv", PostURL: gb0TribuneURL + "/gb0/post", PostField: "message"}
	}
	return g
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

func (g *gb3) handleList(w http.ResponseWriter, r *http.Request) {
	tribunes := make([]string, 0, len(g.tribunes))
	for t := range g.tribunes {
		tribunes = append(tribunes, t)
	}
	sort.Strings(tribunes)
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	err := encoder.Encode(tribunes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	terms := r.FormValue("terms")
	results, err := g.indexer.Search(terms)
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
	g.pollTribunes()
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
	g := newGb3()
	go g.forwardLoop()
	go g.pollLoop()
	RegisterLinuxfrAPI()
	http.HandleFunc("/gb2c/list", func(w http.ResponseWriter, r *http.Request) {
		g.handleList(w, r)
	})
	http.HandleFunc("/gb2c/poll", func(w http.ResponseWriter, r *http.Request) {
		g.handlePoll(w, r)
	})
	http.HandleFunc("/gb2c/post", func(w http.ResponseWriter, r *http.Request) {
		g.handlePost(w, r)
	})
	http.HandleFunc("/gb2c/search", func(w http.ResponseWriter, r *http.Request) {
		g.handleSearch(w, r)
	})
	http.HandleFunc("/gb2c/totoz/search", TotozSearch)
	http.HandleFunc("/gb2c/emoji/search", EmojiSearch)

	gc2Path := os.Getenv("GB2C_GC2_FROM_PATH")
	if len(gc2Path) == 0 {
		gc2Path = "../../gc2"
	}
	gc2PathInfo, err := os.Stat(gc2Path)
	if nil == err && gc2PathInfo.IsDir() {
		http.Handle("/", http.FileServer(http.Dir(gc2Path)))
	}

	listenAddress := os.Getenv("GB2C_LISTEN")
	if len(listenAddress) == 0 {
		listenAddress = ":16666"
	}
	log.Printf("Listen to %s\n", listenAddress)
	err = http.ListenAndServe(listenAddress, nil)
	if nil != err {
		log.Fatal(err)
	}
}
