package main

import (
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
	"time"
)

type AnalyticsHandler struct {
	repo IRepository
}

func NewAnalyticsHandler(repo IRepository) *AnalyticsHandler {
	return &AnalyticsHandler{repo: repo}
}

type CreateLogRequest struct {
	TrackID    int     `json:"track_id"`
	AmountPaid float64 `json:"amount_paid"`
}

type UpdatePriceRequest struct {
	NewPrice float64 `json:"new_price"`
}

type TopTrackStat struct {
	Title string `json:"title"`
	Count int    `json:"count"`
}

func (h *AnalyticsHandler) HandleLogPlayback(w http.ResponseWriter, r *http.Request) {
	var req CreateLogRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	_, err := h.repo.GetTrackByID(req.TrackID)
	if err != nil {
		http.Error(w, "Track not found", http.StatusNotFound)
		return
	}

	log := PlaybackLog{
		TrackID:    req.TrackID,
		AmountPaid: req.AmountPaid,
		PlayedAt:   time.Now(),
	}
	h.repo.CreateLog(log)

	w.WriteHeader(http.StatusCreated)
}

func (h *AnalyticsHandler) HandleUpdatePrice(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	trackID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid track ID", http.StatusBadRequest)
		return
	}

	var req UpdatePriceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.NewPrice <= 0 {
		http.Error(w, "Price must be greater than 0", http.StatusBadRequest)
		return
	}

	err = h.repo.UpdateTrackPrice(trackID, req.NewPrice)
	if err != nil {
		http.Error(w, "Track not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *AnalyticsHandler) HandleGetTopTracks(w http.ResponseWriter, r *http.Request) {
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(top3)
}
