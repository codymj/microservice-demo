package model

import "testing"

func TestProduct_Validate(t *testing.T) {
	p := &Product{
		Name: "Coffee",
		Price: 1.00,
		SKU: "76f2c6eb-3244-4390-90e5-9509b6691f85",
	}

	err := p.Validate()
	if err != nil {
		t.Fatal(err)
	}
}
