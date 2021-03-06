package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
	"unicode/utf8"
)

//go:embed static
var staticContent embed.FS

var postStore Store

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
	info := r.FormValue("info")
	if len(info) == 0 {
		info = r.UserAgent()
	}
	if !utf8.ValidString(info) {
		http.Error(w, "Invalid utf8 in user agent", http.StatusBadRequest)
		return
	}
	postStore.WritePost(Post{
		Message: stripControlsCharsFromString(fmt.Sprintf("%.65536s", message)),
		Info:    stripControlsCharsFromString(fmt.Sprintf("%.32s", info)),
	})
}

func handleGetTsv(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/tab-separated-values")
	posts := postStore.ReadPosts()
	for _, p := range posts {
		p.PrintTsv(w)
	}
}

func main() {
	postStore = NewStore()
	http.HandleFunc("/gb0/tsv", handleGetTsv)
	http.HandleFunc("/gb0/post", handlePost)

	staticfs, err := fs.Sub(staticContent, "static")
	if nil != err {
		log.Fatal(err)
	}
	http.Handle("/", http.FileServer(http.FS(staticfs)))

	listenAddress := os.Getenv("GB0_LISTEN")
	if len(listenAddress) == 0 {
		listenAddress = ":16667"
	}
	log.Printf("Listen to %s\n", listenAddress)
	err = http.ListenAndServe(listenAddress, nil)
	if nil != err {
		log.Fatal(err)
	}
}
