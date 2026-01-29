package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type CreateLogRequest struct {
	TrackID    int     `json:"track_id"`
	AmountPaid float64 `json:"amount_paid"`
}

type UpdatePriceRequest struct {
	NewPrice float64 `json:"new_price"`
}

type AnalyticsHandler struct {
	s *Service
}

func NewHandler(s *Service) *AnalyticsHandler {
	return &AnalyticsHandler{s: s}
}

func (h *AnalyticsHandler) HandleLogPlayback(w http.ResponseWriter, r *http.Request) {
	var req CreateLogRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := h.s.CreateLog(req.TrackID, req.AmountPaid)
	if errors.Is(err, TrackNotFoundError) {
		http.Error(w, "Track not found", http.StatusNotFound)
		return
	}
	if errors.Is(err, FailedToCreateLog) {
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

	err = h.s.UpdatePrice(trackID, req.NewPrice)
	if errors.Is(err, TrackNotFoundError) {
		http.Error(w, "Track not found", http.StatusNotFound)
		return
	}
	if errors.Is(err, PriceMustBeGreater) {
		http.Error(w, "Price must be greater than 0", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *AnalyticsHandler) HandleGetTopTracks(w http.ResponseWriter, r *http.Request) {
	top3, err := h.s.GetTopTracks()
	if errors.Is(err, FailedToGetStats) {
		http.Error(w, "Failed to get stats", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(top3)
}
