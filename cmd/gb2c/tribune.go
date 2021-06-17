package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

const (
	// NoAuthentification Anonymous post
	NoAuthentification = 0
	// OAuth2Authentification annoying dlfp authentification
	OAuth2Authentification = 1
)

// Tribune parameters
type Tribune struct {
	Name                 string
	BackendURL           string
	PostURL              string
	PostField            string
	InfoField            string
	AuthentificationType int
}

// DlfpToken OAuth2 token
type DlfpToken struct {
	AccessToken string `json:"access_token"`
}

// Post message to tribune
func (tribune *Tribune) Post(inRequest *http.Request) {
	data := url.Values{
		tribune.PostField: []string{inRequest.PostFormValue("message")},
	}
	if len(tribune.InfoField) > 0 {
		data[tribune.InfoField] = []string{inRequest.PostFormValue("info")}
	}
	outRequest, err := http.NewRequest("POST", tribune.PostURL, strings.NewReader(data.Encode()))
	if nil != err {
		log.Println(err)
		return
	}
	outRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if len(tribune.InfoField) == 0 {
		outRequest.Header.Set("User-Agent", inRequest.PostFormValue("info"))
	}
	if tribune.AuthentificationType == OAuth2Authentification {
		outRequest.Header.Set("Authorization", inRequest.Header.Get("Authorization"))
	}
	_, err = gb2cHttpClient.Do(outRequest)
	if nil != err {
		log.Println(err)
	}
}

//Poll Retrieve backend
func (tribune *Tribune) Poll() (posts Posts, err error) {
	resp, err := gb2cHttpClient.Get(tribune.BackendURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	posts, err = tribune.parseTSVBackend(resp.Body)
	if err != nil {
		return nil, err
	}

	sort.Sort(posts)
	return posts, nil
}

func (tribune *Tribune) parseTSVBackend(body io.ReadCloser) (Posts, error) {
	posts := make(Posts, 0, MaxPostsInHistory)
	data, err := ioutil.ReadAll(body)
	if nil != err {
		return nil, err
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		record := strings.Split(line, "\t")
		if len(record) < 5 {
			continue
		}
		id, err := strconv.ParseInt(record[0], 10, 64)
		if nil != err {
			log.Println(err)
			continue
		}
		post := Post{ID: id, Time: record[1], Info: record[2], Login: record[3], Message: record[4], Tribune: tribune.Name}
		posts = append(posts, post)
	}
	return posts, nil
}
