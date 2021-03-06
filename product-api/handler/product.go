package handler

import (
	"context"
	"fmt"
	"github.com/codymj/microservice-demo/product-api/model"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

// ProductHandler is an http.Handler
type ProductHandler struct {
	logger *log.Logger
}

// GetProductHandler creates a products handler with the given logger
func GetProductHandler(logger *log.Logger) *ProductHandler {
	return &ProductHandler{logger}
}

// GetAllProducts returns the products from the data store
func (p *ProductHandler) GetAllProducts(rw http.ResponseWriter, r *http.Request) {
	p.logger.Println("Getting products")

	products := model.GetAllProducts()
	err := products.ToJSON(rw)
	if err != nil {
		p.logger.Println("[ERROR] Unable to encode JSON")
		http.Error(rw, "Unable to encode JSON", http.StatusInternalServerError)
	}
}

// AddProduct creates a product and saves to the data store
func (p *ProductHandler) AddProduct(rw http.ResponseWriter, r *http.Request) {
	p.logger.Println("Adding product")

	product := r.Context().Value(KeyProduct{}).(model.Product)
	model.AddProduct(&product)
}

// UpdateProduct updates a product by id
func (p *ProductHandler) UpdateProduct(rw http.ResponseWriter, r *http.Request) {
	// Extract id from URI
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		p.logger.Println("[ERROR] Unable to convert path id")
		http.Error(rw, "Unable to convert path id", http.StatusBadRequest)
	}

	p.logger.Println("Updating product:", id)

	product := r.Context().Value(KeyProduct{}).(model.Product)
	err = model.UpdateProduct(id, &product)
	if err == model.ErrProductNotFound {
		p.logger.Println(fmt.Sprintf("[ERROR] Product %d not found", id))
		http.Error(rw, fmt.Sprintf("Product %d not found", id), http.StatusNotFound)
		return
	}
	if err != nil {
		p.logger.Println("[ERROR] Unexpected error")
		http.Error(rw, "Unexpected error", http.StatusInternalServerError)
		return
	}
}

type KeyProduct struct {}
func (p *ProductHandler) MWValidateProduct(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// Instantiate a Product
		product := model.Product{}

		// Attempt to unmarshal JSON to product
		err := product.FromJSON(r.Body)
		if err != nil {
			p.logger.Println("[ERROR] Unable to decode JSON")
			http.Error(rw, "Unable to decode JSON", http.StatusBadRequest)
			return
		}

		// Validate the product
		err = product.Validate()
		if err != nil {
			p.logger.Println(fmt.Sprintf("Invalid product format: %s", err))
			http.Error(rw, fmt.Sprintf("Invalid product format: %s", err), http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), KeyProduct{}, product)
		r = r.WithContext(ctx)

		next.ServeHTTP(rw, r)
	})
}
