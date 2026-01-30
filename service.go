package main

import (
	"errors"
	"fmt"
	"time"
)

const topTracks = 3

var TrackNotFoundError = errors.New("track not found")
var FailedToCreateLog = errors.New("failed to create log")
var PriceMustBeGreater = errors.New("price must be greater than 0")
var FailedToGetStats = errors.New("failed to get stats")

type Service struct {
	repo IRepository
}

func NewService(repo IRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateLog(trackID int, amountPaid float64) error {
	_, err := s.repo.GetTrackByID(trackID)
	if err != nil {
		return TrackNotFoundError
	}

	log := PlaybackLog{
		TrackID:    trackID,
		AmountPaid: amountPaid,
		PlayedAt:   time.Now(),
	}

	if err := s.repo.CreateLog(log); err != nil {
		return fmt.Errorf("%w: %v", FailedToCreateLog, err)
	}

	return nil
}

func (s *Service) UpdatePrice(trackID int, newPrice float64) error {
	if newPrice <= 0 {
		return PriceMustBeGreater
	}

	err := s.repo.UpdateTrackPrice(trackID, newPrice)
	if err != nil {
		return TrackNotFoundError
	}

	return nil
}

func (s *Service) GetTopTracks() ([]TopTrackStat, error) {
	top3, err := s.repo.GetTopTracks(topTracks)
	if err != nil {
		return nil, FailedToGetStats
	}

	return top3, nil
}
