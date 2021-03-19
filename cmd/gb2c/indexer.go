package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/blevesearch/bleve/v2"
)

//Indexer index and query tribune's posts
type Indexer interface {
	Index(posts Posts)
	Search(query string) (SearchResults, error)
}

//SearchResults query to index results
type SearchResults struct {
	Posts Posts       `json:"posts"`
	Debug interface{} `json:"debug"`
}

type noopIndexer struct {
}

type bleveIndexer struct {
	index bleve.Index
}

//NewIndexer create indexer
func NewIndexer() Indexer {
	//indexPath := os.Getenv("GB2C_FULLTEXT_SEARCH_INDEX_PATH")
	indexPath, _ := os.MkdirTemp(os.TempDir(), "gb3index")
	indexPath = indexPath + "/prout"
	if len(indexPath) == 0 {
		log.Println("Do not forget to define GB2C_FULLTEXT_SEARCH_INDEX_PATH environment variable if you need posts search feature")
		return &noopIndexer{}
	}
	index, err := openOrCreateIndex(indexPath)
	if nil != err {
		log.Println(err)
		return &noopIndexer{}
	}
	return &bleveIndexer{index: index}
}

func openOrCreateIndex(indexPath string) (bleve.Index, error) {
	_, err := os.Stat(indexPath)
	if errors.Is(err, os.ErrNotExist) {
		mapping := bleve.NewIndexMapping()
		return bleve.New(indexPath, mapping)
	} else {
		return bleve.Open(indexPath)
	}
}

func (i *bleveIndexer) Index(posts Posts) {
	for _, p := range posts {
		id := fmt.Sprintf("%d@%s", p.ID, p.Tribune)
		js, err := json.Marshal(p)
		if nil != err {
			log.Println("Cannot marshall post for index ", p, " :", err)
			continue
		}
		err = i.index.Index(id, p)
		if nil != err {
			log.Println("Cannot index post ", p, " :", err)
			continue
		}
		err = i.index.SetInternal([]byte(id), js)
		if nil != err {
			log.Println("Cannot save post in index ", p, " :", err)
			continue
		}
	}
}

func (i *bleveIndexer) Search(query string) (SearchResults, error) {
	matchQuery := bleve.NewQueryStringQuery(query)
	search := bleve.NewSearchRequest(matchQuery)
	bleveResults, err := i.index.Search(search)
	searchResults := SearchResults{Posts: make(Posts, 0)}
	for _, hit := range bleveResults.Hits {
		js, err := i.index.GetInternal([]byte(hit.ID))
		if nil != err {
			log.Println("Cannot retrieve post from index with id : ", hit.ID)
			continue
		}
		var post Post
		err = json.Unmarshal(js, &post)
		if nil != err {
			log.Println("Cannot unmarshal post from index with id : ", hit.ID)
			continue
		}
		searchResults.Posts = append(searchResults.Posts, post)
	}
	return searchResults, err
}

func (i *noopIndexer) Index(posts Posts) {
}

func (i *noopIndexer) Search(query string) (SearchResults, error) {
	return SearchResults{Posts: make(Posts, 0)}, nil
}
