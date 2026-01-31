package main

import (
	"log/slog"
	"net"
	"net/http"
	"os"

	pb "jukebox-analytic/proto"

	"google.golang.org/grpc"
)

const (
	httpPort = ":8080"
	grpcPort = ":50051"
	dbPath   = "./jukebox.db"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	repo, err := NewSQLiteRepository(dbPath)
	if err != nil {
		slog.Error("failed to initialize database", "error", err)
		os.Exit(1)
	}

	service := NewService(repo)

	go func() {
		lis, err := net.Listen("tcp", grpcPort)
		if err != nil {
			slog.Error("failed to listen grpc", "error", err)
			os.Exit(1)
		}

		grpcServer := grpc.NewServer()
		pb.RegisterAnalyticsServiceServer(grpcServer, NewGRPCServer(service))

		slog.Info("gRPC server starting", "address", grpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			slog.Error("grpc server failed", "error", err)
		}
	}()

	handler := NewHandler(service)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/logs", handler.HandleLogPlayback)
	mux.HandleFunc("GET /api/v1/stats/top", handler.HandleGetTopTracks)
	mux.HandleFunc("PATCH /api/v1/tracks/{id}/price", handler.HandleUpdatePrice)

	slog.Info("HTTP server starting", "address", httpPort)
	if err := http.ListenAndServe(httpPort, mux); err != nil {
		slog.Error("http server failed", "error", err)
		os.Exit(1)
	}
}
