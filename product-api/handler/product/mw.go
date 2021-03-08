package product

import (
	"context"
	"fmt"
	"github.com/codymj/microservice-demo/product-api/model"
	"net/http"
)

func (p *Handler) MWValidateProduct(next http.Handler) http.Handler {
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
