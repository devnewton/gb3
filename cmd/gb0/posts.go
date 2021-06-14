package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"time"
)

var timeLocation *time.Location

//Post a moule post
type Post struct {
	Time    string
	Info    string
	Message string
}

//Posts list
type Posts []Post

func (p Posts) Len() int           { return len(p) }
func (p Posts) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p Posts) Less(i, j int) bool { return p[i].Time < p[j].Time }

func (p *Post) PrintTsv(w io.Writer) {
	fmt.Fprintf(w, "%s\t%s\t%s\t\t%s\n", p.Time, p.Time, p.Info, p.Message)
}

func tsvSplit(data []byte, atEOF bool) (advance int, token []byte, err error) {
	for i := 0; i < len(data); i++ {
		d := data[i]
		if d == '\t' || d == '\n' {
			return i + 1, data[:i], nil
		}
	}
	if !atEOF {
		return 0, nil, nil
	}
	return 0, data, bufio.ErrFinalToken
}

func scanTsv(r io.Reader) (*Post, error) {
	scanner := bufio.NewScanner(r)
	scanner.Split(tsvSplit)
	scanner.Scan() //first field is useless id
	scanner.Scan()
	var p Post
	p.Time = scanner.Text()
	scanner.Scan()
	p.Info = scanner.Text()
	scanner.Scan() //next field is useless login
	scanner.Scan()
	p.Message = scanner.Text()
	return &p, scanner.Err()
}

func init() {
	initTimeLocation()
	initMaxPostsInHistory()
	initMaxDaysInHistory()
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

var maxPostsInHistory = 200

func initMaxPostsInHistory() {
	value, isPresent := os.LookupEnv("GB0_MAX_POSTS_IN_HISTORY")
	if isPresent {
		intValue, err := strconv.Atoi(value)
		if nil != err {
			log.Println("Invalid value for GB0_MAX_POSTS_IN_HISTORY")
		} else {
			maxPostsInHistory = intValue
		}
	}
}

var maxDaysInHistory = 7

func initMaxDaysInHistory() {
	value, isPresent := os.LookupEnv("GB0_MAX_DAYS_IN_HISTORY")
	if isPresent {
		intValue, err := strconv.Atoi(value)
		if nil != err {
			log.Println("Invalid value for GB0_MAX_DAYS_IN_HISTORY")
		} else {
			maxDaysInHistory = intValue
		}
	}
}

//NewStore create posts store
func NewStore() Store {
	storeDirectory := os.Getenv("GB0_STORE_DIRECTORY")
	if len(storeDirectory) == 0 {
		log.Println("GB0_STORE_DIRECTORY is not defined, posts will be stored in temporary directory")
		storeDirectory, _ = os.MkdirTemp("", "gb0")
	}
	fileStore := &FileStore{postsDirectory: path.Join(storeDirectory, "posts"), tmpDirectory: path.Join(storeDirectory, "tmp")}
	err := os.MkdirAll(fileStore.postsDirectory, 0755)
	if nil != err {
		log.Fatalf("Cannot create posts directory, try to set GB0_STORE_DIRECTORY to a writable directory : %s\n", err)
	}
	err = os.MkdirAll(fileStore.tmpDirectory, 0755)
	if nil != err {
		log.Fatalf("Cannot create tmp directory, try to set GB0_STORE_DIRECTORY to a writable directory : %s\n", err)
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
	tmpDirectory   string
}

//ReadPosts Load from file
func (f *FileStore) ReadPosts() Posts {
	posts := make(Posts, 0, maxPostsInHistory)
	now := time.Now().In(timeLocation)
	for day := 0; day <= maxDaysInHistory; day++ {
		dayTime := now.AddDate(0, 0, -day)
		dayDir := f.postsDirectory + "/" + dayTime.Format("2006/01/02")
		posts = readPostsFromDirectory(posts, dayDir)
	}
	sort.Sort(posts)
	return posts
}

func readPostsFromDirectory(posts Posts, directory string) Posts {
	files, err := filepath.Glob(directory + "/*.tsv")
	if nil != err {
		log.Printf("Cannot read posts from %s", directory)
		return make(Posts, 0)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(files)))
	for _, file := range files {
		post, err := readPostFromFile(file)
		if nil == err {
			posts = append(posts, *post)
		} else {
			log.Printf("Cannot read posts from %s : %s", file, err)
		}
		if len(posts) >= maxPostsInHistory {
			break
		}
	}
	return posts
}

func readPostFromFile(file string) (*Post, error) {
	tsvFile, err := os.Open(file)
	if nil != err {
		return nil, err
	}
	defer tsvFile.Close()
	var post *Post
	post, err = scanTsv(tsvFile)
	return post, err
}

const maxWritePostTry = 42

//WritePost Save to file
func (f *FileStore) WritePost(p Post) {
	var err error
	for i := 1; i <= maxWritePostTry; i++ {
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
	tmpFile, err := os.CreateTemp(f.tmpDirectory, p.Time+"*.tsv")
	if nil != err {
		return err
	}
	defer tmpFile.Close()
	defer os.Remove(tmpFile.Name())
	p.PrintTsv(tmpFile)

	postDir := f.postsDirectory + "/" + t.Format("2006/01/02")
	err = os.MkdirAll(postDir, 0755)
	if nil != err {
		log.Printf("Cannot create directory %s : %s\n", postDir, err)
		return err
	}
	return os.Link(tmpFile.Name(), path.Join(postDir, p.Time+".tsv"))
}
