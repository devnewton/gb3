package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"sort"
	"time"
)

var storeFile string
var timeLocation *time.Location
var postsLimit int

//Post a moule post
type Post struct {
	ID      int64  `json:"id"`
	Time    string `json:"time"`
	Info    string `json:"info"`
	Message string `json:"message"`
}

//Posts list
type Posts []Post

func (p Posts) Len() int           { return len(p) }
func (p Posts) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p Posts) Less(i, j int) bool { return p[i].ID < p[j].ID }

func init() {
	flag.StringVar(&storeFile, "store", "", "File where to store posts")
	flag.IntVar(&postsLimit, "limit", 200, "Max number of posts to store")
	tl, err := time.LoadLocation("Europe/Paris")
	if nil != err {
		log.Println(err)
		timeLocation = time.UTC
	} else {
		timeLocation = tl
	}
}

//NewStore create file or memory store
func NewStore() Store {
	if len(storeFile) > 0 {
		return &FileStore{}
	}
	return &MemStore{}
}

func setPostTimeAndID(p *Post) {
	t := time.Now().In(timeLocation)
	p.Time = t.Format("20060102150405")
	p.ID = t.UnixNano()
}

//Store posts
type Store interface {
	ReadPosts() (p Posts)
	WritePost(p Post)
}

//FileStore store posts in file
type FileStore struct {
}

//ReadPosts Load from file
func (f *FileStore) ReadPosts() (p Posts) {
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

//WritePost Save to file
func (f *FileStore) WritePost(p Post) {
	setPostTimeAndID(&p)
	newPosts := append(f.ReadPosts(), p)
	sort.Sort(newPosts)
	if len(newPosts) >= postsLimit {
		newPosts = newPosts[len(newPosts)-postsLimit:]
	}
	jsonFile, err := os.OpenFile(storeFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0660)
	if nil != err {
		log.Println(err)
		return
	}
	defer jsonFile.Close()
	encoder := json.NewEncoder(jsonFile)
	encoder.Encode(newPosts)
}

//MemStore store posts in memory
type MemStore struct {
	posts Posts
}

//ReadPosts from memory
func (m *MemStore) ReadPosts() (p Posts) {
	return m.posts
}

//WritePost to memory
func (m *MemStore) WritePost(p Post) {
	setPostTimeAndID(&p)
	m.posts = append(m.posts, p)
	if len(m.posts) >= postsLimit {
		m.posts = m.posts[len(m.posts)-postsLimit:]
	}
}
