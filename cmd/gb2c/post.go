package main

//Post tribune post
type Post struct {
	ID      int64  `json:"id"`
	Time    string `json:"time"`
	Info    string `json:"info"`
	Login   string `json:"login"`
	Message string `json:"message"`
	Tribune string `json:"tribune"`
}

//Posts post array
type Posts []Post

func (p Posts) Len() int           { return len(p) }
func (p Posts) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p Posts) Less(i, j int) bool { return p[i].ID < p[j].ID }
