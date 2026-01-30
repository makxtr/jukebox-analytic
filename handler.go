package main

import (
	"encoding/json"
	"errors"
	"log/slog"
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

func respondWithError(w http.ResponseWriter, r *http.Request, status int, message string, err error, details slog.Attr) {
	slog.Error(message, "error", err, "method", r.Method, "path", r.URL.Path, details)
	http.Error(w, message, status)
}

func (h *AnalyticsHandler) HandleLogPlayback(w http.ResponseWriter, r *http.Request) {
	var req CreateLogRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, r, http.StatusBadRequest, "Invalid request body", err, slog.Any("request_body", r.Body))
		return
	}

	err := h.s.CreateLog(req.TrackID, req.AmountPaid)
	if err != nil {
		if errors.Is(err, TrackNotFoundError) {
			respondWithError(w, r, http.StatusNotFound, "Track not found", err, slog.Int("track_id", req.TrackID))
		} else {
			respondWithError(w, r, http.StatusInternalServerError, "Failed to create log", err, slog.Int("track_id", req.TrackID))
		}
		return
	}

	slog.Info("playback logged successfully", "track_id", req.TrackID, "amount_paid", req.AmountPaid)
	w.WriteHeader(http.StatusCreated)
}

func (h *AnalyticsHandler) HandleUpdatePrice(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	trackID, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, r, http.StatusBadRequest, "Invalid track ID", err, slog.String("track_id_str", idStr))
		return
	}

	var req UpdatePriceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, r, http.StatusBadRequest, "Invalid request body", err, slog.Any("request_body", r.Body))
		return
	}

	err = h.s.UpdatePrice(trackID, req.NewPrice)
	if err != nil {
		details := slog.Group("details", slog.Int("track_id", trackID), slog.Float64("new_price", req.NewPrice))
		if errors.Is(err, TrackNotFoundError) {
			respondWithError(w, r, http.StatusNotFound, "Track not found", err, details)
		} else if errors.Is(err, PriceMustBeGreater) {
			respondWithError(w, r, http.StatusBadRequest, "Price must be greater than 0", err, details)
		} else {
			respondWithError(w, r, http.StatusInternalServerError, "Failed to update price", err, details)
		}
		return
	}

	slog.Info("price updated successfully", "track_id", trackID, "new_price", req.NewPrice)
	w.WriteHeader(http.StatusOK)
}

func (h *AnalyticsHandler) HandleGetTopTracks(w http.ResponseWriter, r *http.Request) {
	top3, err := h.s.GetTopTracks()
	if err != nil {
		respondWithError(w, r, http.StatusInternalServerError, "Failed to get stats", err, slog.String("details", "no details"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(top3)
}
