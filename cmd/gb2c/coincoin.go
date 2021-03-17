package main

import (
	"fmt"
	"net/http"
)

//Coincoin connected to gb3
type Coincoin struct {
	w                   http.ResponseWriter
	flusher             http.Flusher
	LastPostIDByTribune map[string]int64
}

//NewCoincoin create coincoin
func NewCoincoin(w http.ResponseWriter, flusher http.Flusher) *Coincoin {
	return &Coincoin{w: w, flusher: flusher, LastPostIDByTribune: make(map[string]int64)}
}

//Write write events to coincoin
func (c *Coincoin) Write(post Post, data []byte) (err error) {
	if post.ID > c.LastPostIDByTribune[post.Tribune] {
		c.LastPostIDByTribune[post.Tribune] = post.ID
		_, err = fmt.Fprintf(c.w, "data: %s\n\n", data)
		c.flusher.Flush()
	}
	return err
}
