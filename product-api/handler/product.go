package handler

import (
	"context"
	"github.com/codymj/microservice-demo/product-api/model"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

// Product is an http.Handler
type Product struct {
	logger *log.Logger
}

// GetProductHandler creates a products handler with the given logger
func GetProductHandler(logger *log.Logger) *Product {
	return &Product{logger}
}

// GetAllProducts returns the products from the data store
func (p *Product) GetAllProducts(rw http.ResponseWriter, r *http.Request) {
	p.logger.Println("Getting products")

	products := model.GetAllProducts()
	err := products.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal JSON", http.StatusInternalServerError)
	}
}

// AddProduct creates a product and saves to the data store
func (p *Product) AddProduct(rw http.ResponseWriter, r *http.Request) {
	p.logger.Println("Adding product")

	product := r.Context().Value(KeyProduct{}).(model.Product)
	model.AddProduct(&product)
}

// UpdateProduct updates a product by id
func (p *Product) UpdateProduct(rw http.ResponseWriter, r *http.Request) {
	// Extract URI id
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert path id", http.StatusBadRequest)
	}

	p.logger.Println("Updating product:", id)

	product := r.Context().Value(KeyProduct{}).(model.Product)
	err = model.UpdateProduct(id, &product)
	if err == model.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(rw, "Product not found", http.StatusInternalServerError)
		return
	}
}

type KeyProduct struct {}
func (p *Product) MWValidateProduct(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		product := model.Product{}

		err := product.FromJSON(r.Body)
		if err != nil {
			http.Error(rw, "Unable to unmarshal JSON", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), KeyProduct{}, product)
		r = r.WithContext(ctx)

		next.ServeHTTP(rw, r)
	})
}
