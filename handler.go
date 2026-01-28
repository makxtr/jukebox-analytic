package main

import (
	"encoding/json"
	"net/http"
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

	if err := h.repo.CreateLog(log); err != nil {
		http.Error(w, "Failed to create log", http.StatusInternalServerError)
		return
	}

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
	top3, err := h.repo.GetTopTracks(3)
	if err != nil {
		http.Error(w, "Failed to get stats", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(top3)
}
