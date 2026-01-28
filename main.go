package main

import (
	"fmt"
	"log"
	"net/http"
)

const serverAddr = ":8080"
const db = "./jukebox.db"

func main() {
	repo, err := NewSQLiteRepository(db)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

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
