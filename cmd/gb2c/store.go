package main

import (
	"sort"
	"sync"
)

//Store tribune's posts
type Store interface {
	Save(tribune string, posts Posts) error
	RetrieveLatests(tribune string) (Posts, error)
}

const (
	//LastestLength is the number of latests posts to store
	LastestLength = 200
)

type memStore struct {
	postsMap *sync.Map
}

func (m memStore) Save(tribune string, posts Posts) error {
	newPosts := make(Posts, 0, LastestLength)
	oldPosts, ok := m.postsMap.Load(tribune)
	if ok {
		newPosts = append(newPosts, oldPosts.(Posts)...)
	}
	newPosts = append(newPosts, posts...)
	sort.Sort(posts)
	if len(newPosts) > LastestLength {
		newPosts = newPosts[len(newPosts)-LastestLength:]
	}
	m.postsMap.Store(tribune, newPosts)
	return nil
}

func (m memStore) RetrieveLatests(tribune string) (Posts, error) {
	posts, _ := m.postsMap.LoadOrStore(tribune, make(Posts, 0, LastestLength))
	return posts.(Posts), nil
}

//NewStore create store
func NewStore() Store {
	return &memStore{postsMap: new(sync.Map)}
}
