package product

import (
	"github.com/codymj/microservice-demo/product-api/model"
	"net/http"
)

// GetAllProducts returns the products from the data store
func (p *Handler) GetAllProducts(rw http.ResponseWriter, r *http.Request) {
	p.logger.Println("Getting products")

	products := model.GetAllProducts()
	err := products.ToJSON(rw)
	if err != nil {
		p.logger.Println("[ERROR] Unable to encode JSON")
		http.Error(rw, "Unable to encode JSON", http.StatusInternalServerError)
	}
}
