package main

import (
	"context"
	"errors"
	"log/slog"

	pb "jukebox-analytic/proto"
)

type GRPCServer struct {
	pb.UnimplementedAnalyticsServiceServer
	service *Service
}

func NewGRPCServer(s *Service) *GRPCServer {
	return &GRPCServer{service: s}
}

func (s *GRPCServer) LogPlayback(ctx context.Context, req *pb.LogPlaybackRequest) (*pb.Empty, error) {
	err := s.service.CreateLog(int(req.TrackId), req.AmountPaid)
	if err != nil {
		slog.Error("grpc: failed to log playback", "error", err)
		return nil, err
	}
	return &pb.Empty{}, nil
}

func (s *GRPCServer) GetTopTracks(ctx context.Context, req *pb.Empty) (*pb.TopTracksResponse, error) {
	stats, err := s.service.GetTopTracks()
	if err != nil {
		slog.Error("grpc: failed to get top tracks", "error", err)
		return nil, err
	}

	var pbStats []*pb.TopTrack
	for _, stat := range stats {
		pbStats = append(pbStats, &pb.TopTrack{
			Title: stat.Title,
			Count: int32(stat.Count),
		})
	}

	return &pb.TopTracksResponse{Tracks: pbStats}, nil
}

func (s *GRPCServer) UpdatePrice(ctx context.Context, req *pb.UpdatePriceRequest) (*pb.Empty, error) {
	err := s.service.UpdatePrice(int(req.TrackId), req.NewPrice)
	if err != nil {
		if errors.Is(err, TrackNotFoundError) {
			slog.Warn("grpc: track not found", "track_id", req.TrackId)
		}
		return nil, err
	}
	return &pb.Empty{}, nil
}
