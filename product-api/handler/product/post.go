package product

import (
	"github.com/codymj/microservice-demo/product-api/model"
	"net/http"
)

// AddProduct creates a product and saves to the data store
func (p *Handler) AddProduct(rw http.ResponseWriter, r *http.Request) {
	p.logger.Println("Adding product")

	product := r.Context().Value(KeyProduct{}).(model.Product)
	model.AddProduct(&product)
}
