package model

import "testing"

func TestProduct_Validate(t *testing.T) {
	p := &Product{
		Name: "Coffee",
		Price: 1.00,
		SKU: "aaa-bbbb-ccc",
	}

	err := p.Validate()
	if err != nil {
		t.Fatal(err)
	}
}
