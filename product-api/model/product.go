package model

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// Product data model
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	SKU         string  `json:"sku"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

// Products is a collection of Product
type Products []*Product

// FromJSON decodes a JSON object into a Product
func (p *Product) FromJSON(reader io.Reader) error {
	decoder := json.NewDecoder(reader)
	return decoder.Decode(p)
}

// ToJSON encodes the contents of the Product to JSON
func (p *Products) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

// GetProducts returns a list of products
func GetAllProducts() Products {
	return products
}

// AddProduct
func AddProduct(p *Product) {
	p.ID = getNextID()
	products = append(products, p)
}

// UpdateProduct
func UpdateProduct(id int, p *Product) error {
	_, i, err := findProductById(id)
	if err != nil {
		return err
	}

	p.ID = id;
	products[i] = p
	return nil
}

// findProduct()
func findProductById(id int) (*Product, int, error) {
	for i, p := range products {
		if p.ID == id {
			return p, i, nil
		}
	}
	return nil, -1, ErrProductNotFound
}

var ErrProductNotFound = fmt.Errorf("product not found")

// getNextID() returns next ID for Products data store
func getNextID() int {
	length := products[len(products) - 1]
	return length.ID + 1
}

// productList is a hard coded list of products for this
var products = []*Product{
	{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       2.45,
		SKU:         "abc323",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
	{
		ID:          2,
		Name:        "Espresso",
		Description: "Short and strong coffee without milk",
		Price:       1.99,
		SKU:         "fjd34",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
}
