package main

import "fmt"

type IRepository interface {
	TrackRepository
	PlaybackLogRepository
}

type inMemoryRepository struct {
	tracks map[int]*Track
	logs   []PlaybackLog
}

func (r *inMemoryRepository) GetTrackByID(id int) (*Track, error) {
	track, ok := r.tracks[id]
	if !ok {
		return nil, fmt.Errorf("track with id %d not found", id)
	}
	return track, nil
}

func (r *inMemoryRepository) UpdateTrackPrice(id int, newPrice float64) error {
	track, err := r.GetTrackByID(id)
	if err != nil {
		return err
	}
	track.Price = newPrice
	return nil
}

func (r *inMemoryRepository) CreateLog(log PlaybackLog) {
	r.logs = append(r.logs, log)
}

func (r *inMemoryRepository) GetAllLogs() []PlaybackLog {
	return r.logs
}

func NewInMemoryRepository() IRepository {
	tracks := map[int]*Track{
		1: {ID: 1, Title: "Dirty Diana", Artist: "Michael Jackson", Price: 1.25},
		2: {ID: 2, Title: "Comfortably Numb", Artist: "Pink Floyd", Price: 1.50},
		3: {ID: 3, Title: "Space Oddity", Artist: "David Bowie", Price: 1.00},
	}
	return &inMemoryRepository{
		tracks: tracks,
		logs:   []PlaybackLog{},
	}
}
