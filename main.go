package main

import (
	"log"
	"net/http"
	"os"

	"github.com/egorgur/go-microsevices-practice/handlers"
)

func main() {
	l := log.New(os.Stdout, "go-server", log.LstdFlags)
	hh := handlers.NewHello(l)

	sm := http.NewServeMux()
	sm.Handle("/", hh)


	http.ListenAndServe(":9090", sm)
}
