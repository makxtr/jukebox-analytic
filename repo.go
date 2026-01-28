package main

import (
	"fmt"
	"sort"
)

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

func (r *inMemoryRepository) CreateLog(log PlaybackLog) error {
	r.logs = append(r.logs, log)
	return nil
}

func (r *inMemoryRepository) GetAllLogs() []PlaybackLog {
	return r.logs
}

func (r *inMemoryRepository) GetTopTracks(limit int) ([]TopTrackStat, error) {
	counts := make(map[int]int)
	for _, log := range r.logs {
		counts[log.TrackID]++
	}

	var stats []TopTrackStat
	for trackID, count := range counts {
		track, err := r.GetTrackByID(trackID)
		if err == nil {
			stats = append(stats, TopTrackStat{Title: track.Title, Count: count})
		}
	}

	sort.Slice(stats, func(i, j int) bool {
		return stats[i].Count > stats[j].Count
	})

	if len(stats) > limit {
		stats = stats[:limit]
	}
	return stats, nil
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
