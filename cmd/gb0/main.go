package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"
)

var listenAddress string
var verboseMode bool

var writePostChan = make(chan Post)
var readPostsChan = make(chan chan Posts)
var postStore Store

func init() {
	flag.StringVar(&listenAddress, "listen", ":16667", "TCP address to listen on")
	flag.BoolVar(&verboseMode, "verbose", false, "Verbose logging")
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
	writePostChan <- Post{
		Message: stripControlsCharsFromString(fmt.Sprintf("%.65536s", message)),
		Info:    stripControlsCharsFromString(fmt.Sprintf("%.32s", info)),
	}
}

func handleGetTsv(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/tab-separated-values")
	replyChan := make(chan Posts)
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

func writeLoop() {
	for {
		select {
		case p := <-writePostChan:
			postStore.WritePost(p)
		case r := <-readPostsChan:
			r <- postStore.ReadPosts()
		}

	}
}

func main() {
	flag.Parse()
	postStore = NewStore()
	http.HandleFunc("/tsv", handleGetTsv)
	http.HandleFunc("/post", handlePost)
	log.Printf("Listen to %s\n", listenAddress)
	go writeLoop()
	err := http.ListenAndServe(listenAddress, nil)
	if nil != err {
		log.Fatal(err)
	}
}