package tribune

import (
	"fmt"
	"net/http"
)

//CoincoinMessage sent to coincoin
type CoincoinMessage struct {
	Post
	Data []byte
}

//Coincoin connected to gb3
type Coincoin struct {
	Send                chan CoincoinMessage
	w                   http.ResponseWriter
	flusher             http.Flusher
	LastPostIDByTribune map[string]int64
}

//NewCoincoin create coincoin
func NewCoincoin(w http.ResponseWriter, flusher http.Flusher) *Coincoin {
	return &Coincoin{w: w, flusher: flusher, Send: make(chan CoincoinMessage, 8), LastPostIDByTribune: make(map[string]int64)}
}

//WriteLoop write events to coincoin
func (c *Coincoin) WriteLoop() {
	closeEvt := c.w.(http.CloseNotifier).CloseNotify()
	for {
		select {
		case <-closeEvt:
			return
		case msg := <-c.Send:
			if msg.ID > c.LastPostIDByTribune[msg.Tribune] {
				c.LastPostIDByTribune[msg.Tribune] = msg.ID
				fmt.Fprintf(c.w, "data: %s\n\n", msg.Data)
			}
		default:
			c.flusher.Flush()
		}
	}
}
