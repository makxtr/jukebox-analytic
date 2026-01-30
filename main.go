package main

import (
	"log/slog"
	"net/http"
	"os"
)

const serverAddr = ":8080"
const db = "./jukebox.db"

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	slog.Info("starting application", "db_path", db)

	repo, err := NewSQLiteRepository(db)
	if err != nil {
		slog.Error("failed to initialize database", "error", err)
		os.Exit(1)
	}

	service := NewService(repo)
	handler := NewHandler(service)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/v1/logs", handler.HandleLogPlayback)
	mux.HandleFunc("GET /api/v1/stats/top", handler.HandleGetTopTracks)
	mux.HandleFunc("PATCH /api/v1/tracks/{id}/price", handler.HandleUpdatePrice)

	slog.Info("server starting", "address", serverAddr)
	if err := http.ListenAndServe(serverAddr, mux); err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}
