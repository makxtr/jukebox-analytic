package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type mockRepository struct {
	tracks                 map[int]*Track
	logs                   []PlaybackLog
	createLogCalled        bool
	updatePriceCalled      bool
	updatePriceArgID       int
	updatePriceArgNewPrice float64
}

func (m *mockRepository) GetTrackByID(id int) (*Track, error) {
	track, ok := m.tracks[id]
	if !ok {
		return nil, fmt.Errorf("track with id %d not found", id)
	}
	return track, nil
}

func (m *mockRepository) UpdateTrackPrice(id int, newPrice float64) error {
	m.updatePriceCalled = true
	m.updatePriceArgID = id
	m.updatePriceArgNewPrice = newPrice

	track, ok := m.tracks[id]
	if !ok {
		return fmt.Errorf("track not found")
	}
	track.Price = newPrice
	return nil
}

func (m *mockRepository) CreateLog(log PlaybackLog) {
	m.createLogCalled = true
	m.logs = append(m.logs, log)
}

func (m *mockRepository) GetAllLogs() []PlaybackLog {
	return m.logs
}

func TestHandleLogPlayback(t *testing.T) {
	t.Run("successful logging", func(t *testing.T) {
		mockRepo := &mockRepository{
			tracks: map[int]*Track{1: {ID: 1, Title: "Test Song"}},
		}
		handler := NewAnalyticsHandler(mockRepo)

		reqBody := []byte(`{"track_id": 1, "amount_paid": 1.25}`)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/logs", bytes.NewBuffer(reqBody))
		w := httptest.NewRecorder()

		handler.HandleLogPlayback(w, req)

		resp := w.Result()
		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Expected status 201 Created, got %d", resp.StatusCode)
		}
		if !mockRepo.createLogCalled {
			t.Error("CreateLog was not called")
		}
	})

	t.Run("track not found", func(t *testing.T) {
		mockRepo := &mockRepository{tracks: map[int]*Track{}}
		handler := NewAnalyticsHandler(mockRepo)

		reqBody := []byte(`{"track_id": 99, "amount_paid": 1.25}`)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/logs", bytes.NewBuffer(reqBody))
		w := httptest.NewRecorder()

		handler.HandleLogPlayback(w, req)

		resp := w.Result()
		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404 Not Found, got %d", resp.StatusCode)
		}
	})
}

func TestHandleUpdatePrice(t *testing.T) {
	t.Run("successful price update", func(t *testing.T) {
		mockRepo := &mockRepository{
			tracks: map[int]*Track{1: {ID: 1, Title: "Test Song", Price: 1.00}},
		}
		handler := NewAnalyticsHandler(mockRepo)

		reqBody := []byte(`{"new_price": 1.50}`)
		req := httptest.NewRequest(http.MethodPatch, "/api/v1/tracks/1/price", bytes.NewBuffer(reqBody))
		req.SetPathValue("id", "1")

		w := httptest.NewRecorder()

		handler.HandleUpdatePrice(w, req)

		resp := w.Result()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200 OK, got %d", resp.StatusCode)
		}
		if mockRepo.updatePriceArgNewPrice != 1.50 {
			t.Errorf("Expected price update to 1.50, got %.2f", mockRepo.updatePriceArgNewPrice)
		}
	})

	t.Run("invalid price", func(t *testing.T) {
		mockRepo := &mockRepository{
			tracks: map[int]*Track{1: {ID: 1, Title: "Test Song"}},
		}
		handler := NewAnalyticsHandler(mockRepo)

		reqBody := []byte(`{"new_price": -5.00}`)
		req := httptest.NewRequest(http.MethodPatch, "/api/v1/tracks/1/price", bytes.NewBuffer(reqBody))
		req.SetPathValue("id", "1")
		w := httptest.NewRecorder()

		handler.HandleUpdatePrice(w, req)

		resp := w.Result()
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400 Bad Request, got %d", resp.StatusCode)
		}
	})
}

func TestHandleGetTopTracks(t *testing.T) {
	mockRepo := &mockRepository{
		tracks: map[int]*Track{
			1: {ID: 1, Title: "Song A"},
			2: {ID: 2, Title: "Song B"},
		},
		logs: []PlaybackLog{
			{TrackID: 1}, {TrackID: 1}, // Song A: 2
			{TrackID: 2}, // Song B: 1
		},
	}
	handler := NewAnalyticsHandler(mockRepo)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/stats/top", nil)
	w := httptest.NewRecorder()

	handler.HandleGetTopTracks(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %d", resp.StatusCode)
	}

	var result []TopTrackStat
	json.NewDecoder(resp.Body).Decode(&result)

	expected := []TopTrackStat{
		{Title: "Song A", Count: 2},
		{Title: "Song B", Count: 1},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Result mismatch.\nExpected: %+v\nGot:      %+v", expected, result)
	}
}
