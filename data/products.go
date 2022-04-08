package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
)

type Product struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float32 `json:"price" validate:"gt=0,required"`
	SKU         string  `json:"sku" validate:"required,sku"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

func (p *Product) FromJSON(r io.Reader) error {
	decoder := json.NewDecoder(r)
	return decoder.Decode(p)
}
func skuValidation(fl validator.FieldLevel) bool {
	rg := regexp.MustCompile("[a-z]+-[a-z]+-[a-z]")
	matches := rg.FindAllString(fl.Field().String(), -1)
	return len(matches) == 1
}

func (p *Product) Validator() error {
	v := validator.New()
	//custom validator
	v.RegisterValidation("sku", skuValidation)
	return v.Struct(p)
}

type Products []*Product

var ErrDeleteError = fmt.Errorf("unable to delete")

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
func UpdateProduct(id uint, prod *Product) error {
	for _, p := range productList {
		if p.ID == id {
			p.ID = prod.ID
			p.Name = prod.Name
			p.CreatedOn = prod.CreatedOn
			p.DeletedOn = prod.DeletedOn
			p.Description = prod.Description
			p.Price = prod.Price
			p.SKU = prod.SKU
			return nil
		}
	}
	return errors.New("unable to update product")
}
func DeleteProduct(id uint) error {
	for i, pr := range productList {
		if pr.ID == id {
			productList[i] = productList[len(productList)-1]
			productList[len(productList)-1] = nil
			productList = productList[:len(productList)-1]
			break
		}
	}
	return nil
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
