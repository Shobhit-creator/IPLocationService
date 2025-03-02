package main

import (
	"log"
	"net/http"

	"github.com/shobhit-Creator/IPLocationService/internal/handlers"
	"github.com/shobhit-Creator/IPLocationService/internal/middleware"
	"github.com/shobhit-Creator/IPLocationService/internal/service"
	"github.com/shobhit-Creator/IPLocationService/internal/workerpool"
)
func main() {
	initConfig()
	service.InitProviders()
	middleware.InitRateLimiter()

	wp := workerpool.NewWorkerPool(5, 100)
	wp.Start()

	// Start the HTTP server
	http.HandleFunc("/api", middleware.RateLimiter(handlers.LocationHandler(wp)))

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}