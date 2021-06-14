package main

import (
	"log"
	"os"
	"sort"
	"strconv"
	"sync"
)

//Store tribune's posts
type Store interface {
	Save(tribune string, posts Posts) error
	RetrieveLatests(tribune string) (Posts, error)
}

//MaxPostsInHistory is the number of latests posts to store
var MaxPostsInHistory = 1000

func init() {
	value, isPresent := os.LookupEnv("GB2C_MAX_POSTS_IN_HISTORY")
	if isPresent {
		intValue, err := strconv.Atoi(value)
		if nil != err {
			log.Println("Invalid value for GB2C_MAX_POSTS_IN_HISTORY")
		} else {
			MaxPostsInHistory = intValue
		}
	}
}

type memStore struct {
	postsMap *sync.Map
}

func constainsById(posts Posts, id int64) bool {
	for _, p := range posts {
		if p.ID == id {
			return true
		}
	}
	return false
}

func appendUniquePosts(posts Posts, elements ...Post) Posts {
	for _, e := range elements {
		if !constainsById(posts, e.ID) {
			posts = append(posts, e)
		}
	}
	return posts
}

func (m memStore) Save(tribune string, posts Posts) error {
	newPosts := make(Posts, 0, MaxPostsInHistory)
	oldPosts, ok := m.postsMap.Load(tribune)
	if ok {
		newPosts = append(newPosts, oldPosts.(Posts)...)
		newPosts = appendUniquePosts(newPosts, posts...)
	} else {
		newPosts = append(newPosts, posts...)
	}
	sort.Sort(newPosts)
	if len(newPosts) > MaxPostsInHistory {
		newPosts = newPosts[len(newPosts)-MaxPostsInHistory:]
	}
	m.postsMap.Store(tribune, newPosts)
	return nil
}

func (m memStore) RetrieveLatests(tribune string) (Posts, error) {
	posts, _ := m.postsMap.LoadOrStore(tribune, make(Posts, 0, MaxPostsInHistory))
	return posts.(Posts), nil
}

//NewStore create store
func NewStore() Store {
	return &memStore{postsMap: new(sync.Map)}
}
