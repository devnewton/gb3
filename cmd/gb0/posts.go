package main

import (
	"encoding/json"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"
)

var timeLocation *time.Location

//Post a moule post
type Post struct {
	Time    string `json:"time"`
	Info    string `json:"info"`
	Message string `json:"message"`
}

//Posts list
type Posts []Post

func (p Posts) Len() int           { return len(p) }
func (p Posts) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p Posts) Less(i, j int) bool { return p[i].Time < p[j].Time }

func init() {
	initTimeLocation()
}

func initTimeLocation() {
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
	storeDirectory := os.Getenv("GB0_STORE_DIRECTORY")
	if len(storeDirectory) > 0 {
		log.Fatal("GB0_STORE_DIRECTORY is not defined")
	}
	fileStore := &FileStore{postsDirectory: path.Join(storeDirectory, "posts")}
	err := os.MkdirAll(fileStore.postsDirectory, 0644)
	if nil != err {
		log.Fatalf("Cannot create file store, try to set GB0_STORE_DIRECTORY to a writable directory : %s\n", err)
	}
	return fileStore
}

//Store posts
type Store interface {
	ReadPosts() (p Posts)
	WritePost(p Post)
}

//FileStore store posts in file
type FileStore struct {
	postsDirectory string
}

//ReadPosts Load from file
func (f *FileStore) ReadPosts() Posts {
	t := time.Now().In(timeLocation)
	todayDir := t.Format("2006/01/02")
	todayPosts := readPostsFromDirectory(todayDir)
	return todayPosts
}

func readPostsFromDirectory(directory string) (posts Posts) {
	files, err := filepath.Glob(directory + "/*.json")
	if nil != err {
		log.Printf("Cannot read posts from %s", directory)
		return posts
	}
	for _, file := range files {
		post, err := readPostFromFile(file)
		if nil == err {
			posts = append(posts, *post)
		} else {
			log.Printf("Cannot read posts from %s", file)
		}
	}
	return posts
}

func readPostFromFile(file string) (*Post, error) {
	jsonFile, err := os.Open(file)
	if nil != err {
		return nil, err
	}
	defer jsonFile.Close()
	decoder := json.NewDecoder(jsonFile)
	var post Post
	err = decoder.Decode(&post)
	return &post, err
}

//WritePost Save to file
func (f *FileStore) WritePost(p Post) {
	var err error
	for i := 1; i <= 42; i++ {
		err = f.tryToWritePost(p)
		if nil == err {
			return
		}
		time.Sleep(time.Duration(i) * time.Second)
	}
	log.Printf("Cannot write post %s : %s\n", p.Time, err)
}

func (f *FileStore) tryToWritePost(p Post) error {
	t := time.Now().In(timeLocation)
	p.Time = t.Format("20060102150405")

	postDir := f.postsDirectory + "/" + t.Format("2006/01/02")
	err := os.MkdirAll(postDir, 0644)
	if nil != err {
		log.Printf("Cannot create directory %s : %s\n", postDir, err)
		return err
	}

	tmpFile, err := os.CreateTemp(f.postsDirectory, p.Time+"*.tmp")
	if nil != err {
		return err
	}
	defer tmpFile.Close()
	defer os.Remove(tmpFile.Name())
	encoder := json.NewEncoder(tmpFile)
	encoder.Encode(p)
	return os.Link(tmpFile.Name(), p.Time+".json")
}
