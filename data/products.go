package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"
)

type Product struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	SKU         string  `json:"sku"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

func (p *Product) FromJSON(r io.Reader) error {
	decoder := json.NewDecoder(r)
	return decoder.Decode(p)
}

type Products []*Product

var DeleteError = fmt.Errorf("Unable to Delete")

func (p *Products) ToJSON(w io.Writer) error {
	encoder := json.NewEncoder(w)
	return encoder.Encode(p)
}

func GetProducts() Products {
	return productList
}
func AddProduct(p *Product) {
	p.ID = GetID()
	productList = append(productList, p)
}
func DeleteProduct(id int) error {
	prod, err := FindById(id)
	if err != nil {
		return DeleteError
	}
	for i, pr := range productList {
		if pr.ID == prod.ID {
			productList[i] = productList[len(productList)-1]
			productList[len(productList)-1] = nil
			productList = productList[:len(productList)-1]
			break
		}
	}
	return nil
}
func FindById(id int) (*Product, error) {
	for _, prod := range productList {
		if prod.ID == uint(id) {
			return prod, nil
		}
	}
	return &Product{}, errors.New("no match found")
}
func GetID() uint {
	last := productList[len(productList)-1].ID
	return last + 1
}

var productList = Products{
	&Product{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       2.45,
		SKU:         "abc323",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
	&Product{
		ID:          2,
		Name:        "Espresso",
		Description: "short and strong coffee without milk",
		Price:       1.99,
		SKU:         "fjh56u",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
}
