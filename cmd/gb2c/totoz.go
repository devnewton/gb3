package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"net/http"
	"net/url"
)

var totozServer string

//Totozes list
type Totozes struct {
	XMLName xml.Name `xml:"totozes" json:"-"`
	Totozes []Totoz  `xml:"totoz" json:"totozes"`
}

//Totoz description
type Totoz struct {
	XMLName xml.Name `xml:"totoz" json:"-"`
	Name    string   `xml:"name" json:"name"`
	Image   string   `json:"image"`
}

func init() {
	flag.StringVar(&totozServer, "totoz-server", "https://nsfw.totoz.eu/", "Totoz server")
}

//TotozSearch handler
func TotozSearch(w http.ResponseWriter, r *http.Request) {
	searchURL := fmt.Sprintf("https://nsfw.totoz.eu/search.xml?terms=%s", url.QueryEscape(r.FormValue("terms")))
	totozSearchResults, err := http.Get(searchURL)
	if nil != err {
		http.Error(w, fmt.Sprintf("Totoz server error: %s. Blame fork.", err), 500)
	}
	defer totozSearchResults.Body.Close()
	decoder := xml.NewDecoder(totozSearchResults.Body)
	var totozes Totozes
	err = decoder.Decode(&totozes)
	if nil != err {
		http.Error(w, fmt.Sprintf("Cannot decode totoz search results: %s. Blame fork.", err), 500)
	}
	for t := range totozes.Totozes {
		u := url.URL{
			Scheme: "https",
			Host:   "nsfw.totoz.eu",
			Path:   fmt.Sprintf("/img/%s", totozes.Totozes[t].Name),
		}
		totozes.Totozes[t].Image = u.String()
	}
	w.Header().Set("Content-Type", "application/json")
	e := json.NewEncoder(w)
	err = e.Encode(&totozes)
	if nil != err {
		http.Error(w, fmt.Sprintf("Cannot encode totoz search results: %s. Blame dave.", err), 500)
	}
}
