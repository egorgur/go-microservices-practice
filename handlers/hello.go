package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Hello struct {
	l *log.Logger
}

func NewHello(l *log.Logger) *Hello {
	return &Hello{l}
}

func (h *Hello) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	h.l.Println("Got a Request")
	d, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, "An error has occurred", http.StatusBadRequest)
		return
	}
	h.l.Printf("Data %s", d)
	fmt.Fprintf(rw, "Response %s", d)
}
