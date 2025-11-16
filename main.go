package main

import (
	"log"
	"net/http"
	"time"

	"knowledge-capsule-api/handlers"
	"knowledge-capsule-api/middleware"
)

func main() {
	mux := http.NewServeMux()

	// Default routes
	mux.HandleFunc("/", handlers.RootHandler)
	mux.HandleFunc("/api", handlers.ApiRootHandler)
	mux.HandleFunc("/health", handlers.HealthHandler)

	// Public routes
	mux.HandleFunc("/api/auth/register", handlers.RegisterHandler)
	mux.HandleFunc("/api/auth/login", handlers.LoginHandler)

	// Protected routes
	mux.Handle("/api/topics", middleware.AuthMiddleware(http.HandlerFunc(handlers.TopicHandler)))
	mux.Handle("/api/capsules", middleware.AuthMiddleware(http.HandlerFunc(handlers.CapsuleHandler)))
	mux.Handle("/api/search", middleware.AuthMiddleware(http.HandlerFunc(handlers.SearchHandler)))

	// Wrap with logger + recover
	handler := middleware.Recover(middleware.Logger(mux))

	server := &http.Server{
		Addr:         ":8080",
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Println("Server running at http://localhost:8080")
	log.Fatal(server.ListenAndServe())
}
