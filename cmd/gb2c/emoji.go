package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hackebrot/turtle"
)

//Emoji description
type Emoji struct {
	Name       string `json:"name"`
	Characters string `json:"characters"`
}

//EmojiSearch handler
func EmojiSearch(w http.ResponseWriter, r *http.Request) {
	turtleEmojis := turtle.Search(r.FormValue("terms"))
	emojis := make([]Emoji, len(turtleEmojis))
	for t, turtleEmoji := range turtleEmojis {
		emojis[t] = Emoji{Name: turtleEmoji.Name, Characters: turtleEmoji.Char}
	}
	w.Header().Set("Content-Type", "application/json")
	e := json.NewEncoder(w)
	err := e.Encode(&emojis)
	if nil != err {
		http.Error(w, fmt.Sprintf("Cannot encode emoji search results: %s. Blame dave.", err), 500)
	}
}
