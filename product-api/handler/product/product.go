package product

import (
	"log"
)

// Handler is an http.Handler
type Handler struct {
	logger *log.Logger
}

// GetHandler creates a products handler with the given logger
func GetHandler(logger *log.Logger) *Handler {
	return &Handler{logger}
}

// KeyProduct a key used for the Product object in the context
type KeyProduct struct {}
