package data

import (
	"context"
	"encoding/json"
	service "goapi/internal/api/service/data"
	"log"
	"net/http"
	"time"
)

// GetDailySummaryHandler retrieves noise data for a specific room and date
// Example: curl -X GET "http://localhost:8080/data/daily/Room1?date=2025-11-07T00:00:00Z" -u admin:password -H "Content-Type: application/json"
func GetByRoomHandler(w http.ResponseWriter, r *http.Request, logger *log.Logger, ds service.DataService) {
	// Get room name from path parameter
	roomName := r.PathValue("room")
	if roomName == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Room name is required"}`))
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	// Get the data for the last 5 weeks from the service
	data, err := ds.GetByRoom(roomName, ctx)
	if err != nil {
		logger.Printf("Could not get weekly summary: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if len(data) == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "No data found for the specified room"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger.Printf("Error encoding weekly summary: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
