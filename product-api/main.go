package main

import (
	"context"
	"github.com/codymj/microservice-demo/product-api/handler/product"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	// Create logger
	logger := log.New(os.Stdout, "products-api ", log.LstdFlags)

	// Create handlers
	productHandler := product.GetHandler(logger)

	// Create new serve mux and register handlers
	serveMux := mux.NewRouter()
	registerProductRouters(serveMux, productHandler)

	// Create and start web server
	s := http.Server{
		Addr:         ":9090",
		Handler:      serveMux,
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	go func() {
		logger.Println("Starting server on port 9090")

		err := s.ListenAndServe()
		if err != nil {
			logger.Printf("Server closed")
			os.Exit(1)
		}
	}()

	// Trap sigterm or interrupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	sig := <-c
	log.Println("Got signal:", sig)

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := s.Shutdown(ctx)
	if err != nil {
		logger.Fatal("Server did not shutdown gracefully")
	}
}

// Sets up sub routers for product
func registerProductRouters(serveMux *mux.Router, productHandler *product.Handler) {
	// GET
	getRouter := serveMux.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/product", productHandler.GetAllProducts)

	// PUT
	putRouter := serveMux.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/product/{id:[0-9]+}", productHandler.UpdateProduct)
	putRouter.Use(productHandler.MWValidateProduct)

	// POST
	postRouter := serveMux.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/product", productHandler.AddProduct)
	postRouter.Use(productHandler.MWValidateProduct)

	// DELETE
	deleteRouter := serveMux.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/product/{id:[0-9]+}", productHandler.DeleteProduct)
}
