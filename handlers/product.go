package handlers

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/cupcake08/grps-golang/data"
	"github.com/gorilla/mux"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET Product")
	lp := data.GetProducts()
	//encode lp into json
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to encode json", http.StatusInternalServerError)
	}
}

func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Product")
	d := r.Context().Value(KeyProduct{}).(*data.Product)
	p.l.Printf("Product: %#v", d)
	data.AddProduct(d)
	rw.Write([]byte("Product added succesfully. :)"))
}

func (p *Products) UpdateProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle PUT Product")
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(rw, "unable to get valid id", http.StatusBadRequest)
		return
	}
	prod := r.Context().Value(KeyProduct{}).(*data.Product)
	data.UpdateProduct(uint(id), prod)
	rw.Write([]byte("Product Updated Successfully. :)"))
}

func (p *Products) DeleteProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle DELETE Product")
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(rw, "unable to get valid id", http.StatusBadRequest)
		return
	}
	if err := data.DeleteProduct(uint(id)); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	rw.Write([]byte("Successfully deleted. :)"))
}

type KeyProduct struct{}

//middleware
func (p *Products) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)
		prod := &data.Product{}
		if err := prod.FromJSON(r.Body); err != nil {
			http.Error(w, "unable to unmarshal json", http.StatusBadRequest)
			return
		}
		//validate the product
		if err := prod.Validator(); err != nil {
			p.l.Println("[ERROR]: ", err)
			http.Error(w, "unable to validate product", http.StatusBadRequest)
			return
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, KeyProduct{}, prod)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
