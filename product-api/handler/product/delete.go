package product

import (
	"fmt"
	"github.com/codymj/microservice-demo/product-api/model"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

// DeleteProduct deletes a product by id
func (p *Handler) DeleteProduct(rw http.ResponseWriter, r *http.Request) {
	// Extract id from URI
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	p.logger.Println("Deleting product:", id)

	err := model.DeleteProduct(id)
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
