package main

import "fmt"

func main() {
	repo := NewInMemoryRepository()

	handler := NewAnalyticsHandler(repo)

	fmt.Println("--- Simulating API Calls ---")

	handler.HandleLogPlayback(1, 1.25)
	handler.HandleLogPlayback(1, 1.25)
	handler.HandleLogPlayback(2, 1.50)
	handler.HandleLogPlayback(3, 1.00)
	handler.HandleLogPlayback(1, 1.25)
	handler.HandleLogPlayback(2, 1.50)
	handler.HandleLogPlayback(4, 1.00)
	fmt.Println()

	topTracks := handler.HandleGetTopTracks()
	fmt.Println("--- Top 3 Tracks ---")
	for _, stat := range topTracks {
		fmt.Printf("- %s: %d plays\n", stat.Title, stat.Count)
	}
	fmt.Println()

	handler.HandleUpdatePrice(1, 0)
	handler.HandleUpdatePrice(1, 1.35)

	track, _ := repo.GetTrackByID(1)
	fmt.Printf("\nFinal price for track %d: %.2f\n", track.ID, track.Price)
}
