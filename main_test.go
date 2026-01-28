package main

import (
	"fmt"
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

		handler.HandleLogPlayback(1, 1.25)

		if !mockRepo.createLogCalled {
			t.Error("CreateLog was not called, but it should have been")
		}
		if len(mockRepo.logs) != 1 {
			t.Errorf("Expected 1 log, but got %d", len(mockRepo.logs))
		}
	})

	t.Run("track not found", func(t *testing.T) {
		mockRepo := &mockRepository{
			tracks: map[int]*Track{},
		}
		handler := NewAnalyticsHandler(mockRepo)

		handler.HandleLogPlayback(99, 1.25)

		if mockRepo.createLogCalled {
			t.Error("CreateLog was called, but it should not have been (track not found)")
		}
	})
}

func TestHandleUpdatePrice(t *testing.T) {
	t.Run("successful price update", func(t *testing.T) {
		mockRepo := &mockRepository{
			tracks: map[int]*Track{1: {ID: 1, Title: "Test Song", Price: 1.00}},
		}
		handler := NewAnalyticsHandler(mockRepo)

		handler.HandleUpdatePrice(1, 1.50)

		if !mockRepo.updatePriceCalled {
			t.Error("UpdateTrackPrice was not called")
		}
		if mockRepo.updatePriceArgID != 1 || mockRepo.updatePriceArgNewPrice != 1.50 {
			t.Errorf("UpdateTrackPrice called with wrong arguments: ID=%d, Price=%.2f",
				mockRepo.updatePriceArgID, mockRepo.updatePriceArgNewPrice)
		}
	})

	t.Run("invalid price (zero)", func(t *testing.T) {
		mockRepo := &mockRepository{
			tracks: map[int]*Track{1: {ID: 1, Title: "Test Song", Price: 1.00}},
		}
		handler := NewAnalyticsHandler(mockRepo)

		handler.HandleUpdatePrice(1, 0)

		if mockRepo.updatePriceCalled {
			t.Error("UpdateTrackPrice was called with an invalid price, but it should not have been")
		}
	})
}

func TestHandleGetTopTracks(t *testing.T) {
	mockRepo := &mockRepository{
		tracks: map[int]*Track{
			1: {ID: 1, Title: "Song A"},
			2: {ID: 2, Title: "Song B"},
			3: {ID: 3, Title: "Song C"},
			4: {ID: 4, Title: "Song D"},
		},
		logs: []PlaybackLog{
			{TrackID: 1}, {TrackID: 1}, {TrackID: 1},
			{TrackID: 2}, {TrackID: 2},
			{TrackID: 3},
			{TrackID: 4},
		},
	}
	handler := NewAnalyticsHandler(mockRepo)

	result := handler.HandleGetTopTracks()

	expected := []TopTrackStat{
		{Title: "Song A", Count: 3},
		{Title: "Song B", Count: 2},
		{Title: "Song C", Count: 1},
	}

	if len(result) != 3 {
		t.Fatalf("Expected 3 tracks in top, but got %d", len(result))
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Result does not match expected.\nExpected: %+v\nGot:      %+v", expected, result)
	}
}
