package data

import (
	"testing"
)

func TestCheckValidation(t *testing.T) {
	p := &Product{
		ID:    3,
		Name:  "americano",
		Price: 1.49,
		SKU:   "abc8-xyz-kjl",
	}
	if err := p.Validator(); err != nil {
		t.Fatal(err)
	}
}
