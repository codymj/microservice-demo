package product

import (
	"fmt"
	"github.com/codymj/microservice-demo/product-api/model"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

// UpdateProduct updates a product by id
func (p *Handler) UpdateProduct(rw http.ResponseWriter, r *http.Request) {
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
