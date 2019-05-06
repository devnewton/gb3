package main

import (
	"flag"
	"log"
	"net/http"
)

var listenAddress string
var verboseMode bool

func init() {
	flag.StringVar(&listenAddress, "listen", ":6667", "TCP address to listen on")
	flag.BoolVar(&verboseMode, "verbose", false, "Verbose logging")
}

type gb0 struct {
}

func newGb0() *gb0 {
	return &gb0{}
}

func (g *gb0) handlePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if nil != err {
		log.Println(err)
		return
	}
	//TODO
}

func (g *gb0) handleGetTsv(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/tab-separated-values")
	//TODO
}

func main() {
	flag.Parse()
	g := newGb0()
	http.HandleFunc("/tsv", func(w http.ResponseWriter, r *http.Request) {
		g.handleGetTsv(w, r)
	})
	http.HandleFunc("/api/post", func(w http.ResponseWriter, r *http.Request) {
		g.handlePost(w, r)
	})
	log.Printf("Listen to %s\n", listenAddress)
	err := http.ListenAndServe(listenAddress, nil)
	if nil != err {
		log.Fatal(err)
	}
}
