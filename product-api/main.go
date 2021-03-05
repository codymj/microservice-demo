package main

import (
	"context"
	"github.com/codymj/microservice-demo/product-api/handler"
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
	productsHandler := handler.GetProductsHandler(logger)

	// Create new serve mux and register handlers
	serveMux := http.NewServeMux()
	serveMux.Handle("/", productsHandler)

	// Create a new server
	s := http.Server{
		Addr:			":9090",
		Handler:		serveMux,
		ErrorLog:		logger,
		ReadTimeout:	5*time.Second,
		WriteTimeout:	10*time.Second,
		IdleTimeout:	120*time.Second,
	}

	// Start the server
	go func() {
		logger.Println("Starting server on port 9090")

		err := s.ListenAndServe()
		if err != nil {
			logger.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	// Trap sigterm or interrupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received
	sig := <-c
	log.Println("Got signal:", sig)

	// Gracefully shutdown the server
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := s.Shutdown(ctx)
	if err != nil {
		logger.Fatal("Server did not shutdown gracefully")
	}
}
