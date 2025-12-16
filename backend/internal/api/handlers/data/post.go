package data

import (
	"context"
	"encoding/json"
	"goapi/internal/api/repository/models"
	service "goapi/internal/api/service/data"
	"log"
	"net/http"
	"time"
)

// PostHandler handles real-time sound data from Arduino
func PostHandler(w http.ResponseWriter, r *http.Request, logger *log.Logger, ds service.DataService) {
	var data models.Data

	// Decode JSON payload
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Invalid request data. Please check your input."}`))
		return
	}

	// Fill timestamp if missing
	if data.MeasureTime == "" {
		data.MeasureTime = time.Now().Format(time.RFC3339)
	}

	// Fill default device ID if missing
	if data.DeviceID == "" {
		data.DeviceID = "arduino_001"
	}

	if data.Threshold == 0 {
		data.Threshold = 70 // default threshold
	}

	logger.Println("Received POST /api/data from Arduino:")
	logger.Printf("%+v\n", data)

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	// Constantly sent data for current sound level card in UI
	// Only has 1 row per device that gets updated
	// So it doesn't need to be periodically cleared
	if !data.IsPeriodic {
		// Save latest_data to repository
		if err := ds.CreateLatest(&data, ctx); err != nil {
			switch err.(type) {
			case service.DataError:
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"error": "` + err.Error() + `"}`))
				return
			default:
				logger.Println("Error updating latest measurement data:", err, data)
				http.Error(w, "Internal server error.", http.StatusInternalServerError)
				return
			}
		}
	} else {
		// Data sent periodically (every 10 minutes)
		// Save data to repository
		if err := ds.Create(&data, ctx); err != nil {
			switch err.(type) {
			case service.DataError:
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"error": "` + err.Error() + `"}`))
				return
			default:
				logger.Println("Error creating data:", err, data)
				http.Error(w, "Internal server error.", http.StatusInternalServerError)
				return
			}
		}
	}

	// Return the created record as JSON
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger.Println("Error encoding data:", err, data)
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}
}
