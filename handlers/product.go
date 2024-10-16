package handlers

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/egorgur/go-microservices-practice/product-api/data"
	"github.com/gorilla/mux"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

// func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
// 	if r.Method == http.MethodGet {
// 		p.l.Println("Handle GET request")
// 		p.getProducts(rw, r)
// 		return
// 	}

// 	if r.Method == http.MethodPost {
// 		p.l.Println("Handle POST request")
// 		p.addProduct(rw, r)
// 		return
// 	}

// 	if r.Method == http.MethodPut {
// 		// expect the id in the URI
// 		p.l.Println("Handle PUT request")
// 		re := regexp.MustCompile(`/([0-9]+)`)
// 		foundMatches := re.FindAllStringSubmatch(r.URL.Path, -1)

// 		if len(foundMatches) != 1 {
// 			http.Error(rw, "Invalid URL", http.StatusBadRequest)
// 			p.l.Printf("Invalid URL parsing %#v", foundMatches)
// 			return
// 		}

// 		if len(foundMatches[0]) != 2 {
// 			http.Error(rw, "Invalid URL", http.StatusBadRequest)
// 			p.l.Printf("Invalid URL parsing %#v", foundMatches)
// 			return
// 		}

// 		idString := foundMatches[0][1]
// 		id, err := strconv.Atoi(idString)
// 		if err != nil {
// 			http.Error(rw, "Failed to parse index from URL", http.StatusBadRequest)
// 			p.l.Fatalf("Invalid URL parsing %#v", idString)
// 			return
// 		}

// 		p.updateProduct(rw, r, id)
// 	}

// 	// catch all not implemented
// 	rw.WriteHeader(http.StatusMethodNotAllowed)
// }

func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	lp := data.GetProducts()
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
		return
	}
}

func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
	prod := r.Context().Value(KeyProduct{}).(data.Product)
	p.l.Printf("Prod %#v", prod)
	data.AddProduct(&prod)
}

func (p *Products) UpdateProduct(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idString := vars["id"]
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(rw, "Failed to parse index from URL", http.StatusBadRequest)
		p.l.Fatalf("Invalid URL parsing %#v", idString)
		return
	}

	prod := r.Context().Value(KeyProduct{}).(data.Product)
	err = data.UpdateProduct(&prod, id)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Unable to found the resource", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(rw, "Unable to update the resource", http.StatusInternalServerError)
		return
	}
}

type KeyProduct struct{}

func (p Products) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		prod := data.Product{}

		err := prod.FromJSON(r.Body)
		if err != nil {
			http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		req := r.WithContext(ctx)
		next.ServeHTTP(rw, req)
	})
}
