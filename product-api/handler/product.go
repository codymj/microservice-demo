package handler

import (
	"github.com/codymj/microservice-demo/product-api/model"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

// Products is an http.Handler
type Products struct {
	logger *log.Logger
}

// NewProducts creates a products handler with the given logger
func GetProductsHandler(logger *log.Logger) *Products {
	return &Products{logger}
}

// ServeHTTP is the entry point for the handler and satisfies the http.Handler
func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// GET
	if r.Method == http.MethodGet {
		p.getAllProducts(rw, r)
		return
	}
	// POST
	if r.Method == http.MethodPost {
		p.addProduct(rw, r)
		return
	}
	// PUT
	if r.Method == http.MethodPut {
		regex := regexp.MustCompile(`/([0-9]+)`)
		group := regex.FindAllStringSubmatch(r.URL.Path, -1)
		if len(group) != 1 || len(group[0]) != 2 {
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}

		idString := group[0][1]
		id, err := strconv.Atoi(idString)
		if err != nil {
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}

		p.updateProduct(id, rw, r)
	}

	rw.WriteHeader(http.StatusMethodNotAllowed)
}

// getAllProducts returns the products from the data store
func (p *Products) getAllProducts(rw http.ResponseWriter, r *http.Request) {
	p.logger.Println("Getting products")

	products := model.GetAllProducts()
	err := products.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal JSON", http.StatusInternalServerError)
	}
}

// addProduct creates a product and saves to the data store
func (p *Products) addProduct(rw http.ResponseWriter, r *http.Request) {
	p.logger.Println("Adding product")

	product := &model.Product{}
	err := product.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal JSON", http.StatusBadRequest)
	}

	model.AddProduct(product)
}

// updateProduct updates a product by ID
func (p *Products) updateProduct(id int, rw http.ResponseWriter, r *http.Request) {
	p.logger.Println("Updating product:", id)

	product := &model.Product{}
	err := product.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal JSON", http.StatusBadRequest)
	}

	err = model.UpdateProduct(id, product)
	if err == model.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(rw, "Product not found", http.StatusInternalServerError)
		return
	}
}
