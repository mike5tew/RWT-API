package main

import (
	"net/http"
)

// Add this to your main.go file or use as reference
func setupHealthEndpoint() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
}

// In your main() function, call setupHealthEndpoint()
