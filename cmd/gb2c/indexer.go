package main

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

//NewIndexer create indexer
func NewIndexer() Indexer {
	return &noopIndexer{}
}

func (i *noopIndexer) Index(posts Posts) {
}

func (i *noopIndexer) Search(query string) (SearchResults, error) {
	return SearchResults{Posts: make(Posts, 0)}, nil
}
