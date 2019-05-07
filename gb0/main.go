package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
	"unicode/utf8"
)

var listenAddress string
var verboseMode bool
var storeFile string
var timeLocation *time.Location
var writePostChan = make(chan post)
var readPostsChan = make(chan chan posts)
var postsLimit int

func init() {
	flag.StringVar(&listenAddress, "listen", ":6667", "TCP address to listen on")
	flag.BoolVar(&verboseMode, "verbose", false, "Verbose logging")
	flag.StringVar(&storeFile, "file", "/tmp/gb0.json", "File to store posts")
	flag.IntVar(&postsLimit, "limit", 200, "Max number of posts to store")
	tl, err := time.LoadLocation("Europe/Paris")
	if nil != err {
		log.Println(err)
		timeLocation = time.UTC
	} else {
		timeLocation = tl
	}
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

func handlePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if nil != err {
		http.Error(w, fmt.Sprintf("Invalid request : %s", err), http.StatusBadRequest)
		return
	}
	message := r.FormValue("message")
	if !utf8.ValidString(message) {
		http.Error(w, "Invalid utf8 in message", http.StatusBadRequest)
		return
	}
	info := r.UserAgent()
	if !utf8.ValidString(info) {
		http.Error(w, "Invalid utf8 in user agent", http.StatusBadRequest)
		return
	}
	writePostChan <- post{
		Message: stripControlsCharsFromString(fmt.Sprintf("%.65536s", message)),
		Info:    stripControlsCharsFromString(fmt.Sprintf("%.32s", info)),
	}
}

func handleGetTsv(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/tab-separated-values")
	replyChan := make(chan posts)
	readPostsChan <- replyChan
	select {
	case postsToSend := <-replyChan:
		for _, p := range postsToSend {
			fmt.Fprintf(w, "%d\t%s\t%s\t%s\t\n", p.ID, p.Time, p.Info, p.Message)
		}
	case <-time.After(30 * time.Second):
		log.Println("Read posts timeout")
	}

}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func writePost(p post) {
	t := time.Now().In(timeLocation)
	p.Time = t.Format("20060102150405")
	p.ID = t.UnixNano()
	newPosts := append(readPosts(), p)
	sort.Sort(newPosts)
	newPosts = newPosts[0:min(postsLimit, len(newPosts))]

	jsonFile, err := os.OpenFile(storeFile, os.O_WRONLY|os.O_CREATE, 0660)
	if nil != err {
		log.Println(err)
		return
	}
	defer jsonFile.Close()
	encoder := json.NewEncoder(jsonFile)
	encoder.Encode(newPosts)
}

func writeLoop() {
	for {
		select {
		case p := <-writePostChan:
			writePost(p)
		case r := <-readPostsChan:
			r <- readPosts()
		}

	}
}

func main() {
	flag.Parse()
	http.HandleFunc("/tsv", handleGetTsv)
	http.HandleFunc("/post", handlePost)
	log.Printf("Listen to %s\n", listenAddress)
	go writeLoop()
	err := http.ListenAndServe(listenAddress, nil)
	if nil != err {
		log.Fatal(err)
	}
}
