package main

import (
	"fmt"
	"sort"
	"time"
)

type AnalyticsHandler struct {
	repo IRepository
}

func NewAnalyticsHandler(repo IRepository) *AnalyticsHandler {
	return &AnalyticsHandler{repo: repo}
}

// [POST] /api/v1/logs
func (h *AnalyticsHandler) HandleLogPlayback(trackID int, amountPaid float64) {
	_, err := h.repo.GetTrackByID(trackID)
	if err != nil {
		fmt.Printf("Error: %v (HTTP 404)\n", err)
		return
	}

	log := PlaybackLog{
		TrackID:    trackID,
		AmountPaid: amountPaid,
		PlayedAt:   time.Now(),
	}
	h.repo.CreateLog(log)

	fmt.Printf("Successfully logged playback for track %d\n", trackID)
}

// [PATCH] /api/v1/tracks/{id}/price
func (h *AnalyticsHandler) HandleUpdatePrice(trackID int, newPrice float64) {
	if newPrice <= 0 {
		fmt.Printf("Error: price must be greater than 0 (HTTP 400)\n")
		return
	}

	err := h.repo.UpdateTrackPrice(trackID, newPrice)
	if err != nil {
		fmt.Printf("Error: %v (HTTP 404)\n", err)
		return
	}

	fmt.Printf("Price for track %d updated successfully\n", trackID)
}

// TopTrackStat
type TopTrackStat struct {
	Title string
	Count int
}

// [GET] /api/v1/stats/top
func (h *AnalyticsHandler) HandleGetTopTracks() []TopTrackStat {
	logs := h.repo.GetAllLogs()
	counts := make(map[int]int)
	for _, log := range logs {
		counts[log.TrackID]++
	}

	var stats []TopTrackStat
	for trackID, count := range counts {
		track, err := h.repo.GetTrackByID(trackID)
		if err == nil {
			stats = append(stats, TopTrackStat{Title: track.Title, Count: count})
		}
	}

	sort.Slice(stats, func(i, j int) bool {
		return stats[i].Count > stats[j].Count
	})

	// TOP-3
	top3 := stats
	if len(top3) > 3 {
		top3 = top3[:3]
	}

	return top3
}
