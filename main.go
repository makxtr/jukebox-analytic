package main

import (
	"fmt"
	"net/http"
)

const serverAddr = ":8080"

func main() {
	repo := NewInMemoryRepository()
	handler := NewAnalyticsHandler(repo)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/v1/logs", handler.HandleLogPlayback)
	mux.HandleFunc("GET /api/v1/stats/top", handler.HandleGetTopTracks)
	mux.HandleFunc("PATCH /api/v1/tracks/{id}/price", handler.HandleUpdatePrice)

	fmt.Printf("Server starting on http://localhost%s\n", serverAddr)
	if err := http.ListenAndServe(serverAddr, mux); err != nil {
		fmt.Printf("Server failed: %v\n", err)
	}
}
