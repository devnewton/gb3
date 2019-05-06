package gb0

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

var listenAddress string
var verboseMode bool
var storeFile string

func init() {
	flag.StringVar(&listenAddress, "listen", ":6667", "TCP address to listen on")
	flag.BoolVar(&verboseMode, "verbose", false, "Verbose logging")
	flag.StringVar(&storeFile, "file", "/tmp/gb0.json", "File to store posts")
}

type post struct {
	ID      int64  `json:"id"`
	Time    string `json:"time"`
	Info    string `json:"info"`
	Message string `json:"message"`
}

type posts []post

func (p posts) Len() int           { return len(p) }
func (p posts) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p posts) Less(i, j int) bool { return p[i].ID < p[j].ID }

type gb0 struct {
}

func newGb0() *gb0 {
	return &gb0{}
}

func readPosts() (p posts) {
	jsonFile, err := os.Open(storeFile)
	if nil != err {
		log.Println(err)
		return p
	}
	defer jsonFile.Close()
	decoder := json.NewDecoder(jsonFile)
	err = decoder.Decode(&p)
	if nil != err {
		return p
	}
	return p
}

func stripControlsCharsFromString(str string) string {
	return strings.Map(func(r rune) rune {
		if r >= 32 && r != 127 {
			return r
		}
		return -1
	}, str)
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
	for _, p := range readPosts() {
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t\n", p.ID, p.Time, p.Info, p.Message)
	}
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
