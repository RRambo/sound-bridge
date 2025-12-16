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
func GetDailySummaryHandler(w http.ResponseWriter, r *http.Request, logger *log.Logger, ds service.DataService) {
	// Get room name from path parameter
	roomName := r.PathValue("room")
	if roomName == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Room name is required"}`))
		return
	}

	// Get date from query parameter
	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Date parameter is required"}`))
		return
	}

	// Parse the date
	date, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Invalid date format. Use RFC3339 format (e.g., 2025-11-07T00:00:00Z)"}`))
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	// Get the daily summary from the service
	data, err := ds.GetDailySummary(roomName, date, ctx)
	if err != nil {
		logger.Printf("Could not get daily summary: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if len(data) == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "No data found for the specified room and date"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger.Printf("Error encoding daily summary: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
