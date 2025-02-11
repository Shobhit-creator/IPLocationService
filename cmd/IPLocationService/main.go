package main

import (
	"net/http"

	"github.com/shobhit-Creator/IPLocationService/internal/handlers"
	"github.com/shobhit-Creator/IPLocationService/internal/middleware"
)
func main() {
	http.HandleFunc("/", middleware.RateLimiter(handlers.LocationHandler))
}