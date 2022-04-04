package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/cupcake08/grps-golang/data"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		p.GetProducts(w, r)
		return
	}

	if r.Method == http.MethodPost {
		p.AddProduct(w, r)
		return
	}

	if r.Method == http.MethodPut {
		//extract id from path
		p.l.Println("extracting id...")
		reg := regexp.MustCompile("/([0-9]+)")
		g := reg.FindAllStringSubmatch(r.URL.Path, -1)

		p.l.Printf("group: %+v", g)

		if len(g) != 1 {
			http.Error(w, "invalid url", http.StatusBadRequest)
			return
		}
		if len(g[0]) != 2 {
			http.Error(w, "invalid url", http.StatusBadRequest)
			return
		}

		idString := g[0][1]
		id, err := strconv.Atoi(idString)
		if err != nil {
			http.Error(w, "invalid url", http.StatusBadRequest)
			return
		}
		p.l.Println("got id", id)
		p.UpdateProduct(id, w, r)
		return
	}
	//catch all other methods
	w.WriteHeader(http.StatusMethodNotAllowed)
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
	d := &data.Product{}

	if err := d.FromJSON(r.Body); err != nil {
		http.Error(rw, "unable to unmarshal data", http.StatusBadRequest)
		return
	}
	p.l.Printf("Product: %#v", d)
	data.AddProduct(d)
	rw.Write([]byte("Product added succesfully. :)"))
}

func (p *Products) UpdateProduct(id int, rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle PUT Product")

	prod, err := data.FindById(id)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}

	d := &data.Product{}
	if err = d.FromJSON(r.Body); err != nil {
		http.Error(rw, "unable to unmarshal json", http.StatusBadRequest)
		return
	}

	prod.ID = d.ID
	prod.Name = d.Name
	prod.Price = d.Price
	prod.Description = d.Description
	prod.CreatedOn = d.CreatedOn
	prod.DeletedOn = d.DeletedOn
	prod.UpdatedOn = d.UpdatedOn

	rw.Write([]byte("Product Updated Successfully. :)"))
}
