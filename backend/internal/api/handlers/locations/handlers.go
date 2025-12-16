package locations

import (
	"context"
	"encoding/json"
	"goapi/internal/api/repository/models"
	service "goapi/internal/api/service/data"
	"log"
	"net/http"
	"strconv"
	"time"
)

func GetLocationsHandler(w http.ResponseWriter, r *http.Request, logger *log.Logger, svc service.LocationService) {
	locations, err := svc.GetAllLocations()
	if err != nil {
		logger.Println("Error getting locations:", err)
		http.Error(w, `{"error": "Failed to get locations"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"locations": locations,
	})
}

type LocationResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func GetChosenLocationHandler(w http.ResponseWriter, r *http.Request, logger *log.Logger, svc service.LocationService) {
	w.Header().Set("Content-Type", "application/json")
	location, err := svc.GetChosenLocation()
	if err != nil {
		logger.Println("Error getting chosen location:", err)
		http.Error(w, `{"error": "Failed to get chosen location"}`, http.StatusInternalServerError)
		return
	}
	if location == nil {
		// * This is a User Error, response in JSON and with a 404 status code
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "Resource not found."}`))
		return
	}

	resp := LocationResponse{
		Message: "Location retrieved",
		Data:    location,
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode response.", http.StatusInternalServerError)
	}
}

func CreateLocationHandler(w http.ResponseWriter, r *http.Request, logger *log.Logger, svc service.LocationService) {
	var location models.Location
	if err := json.NewDecoder(r.Body).Decode(&location); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	if err := svc.CreateLocation(&location); err != nil {
		logger.Println("Error creating location:", err)
		http.Error(w, `{"error": "Failed to create location"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(location)
}

// UpdateThresholdHandler and SetChosenLocationHandler combined
func UpdateLocationHandler(w http.ResponseWriter, r *http.Request, logger *log.Logger, svc service.LocationService) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error": "Invalid location ID"}`, http.StatusBadRequest)
		return
	}

	// Get threshold from query parameter
	thresholdStr := r.URL.Query().Get("newThreshold")
	if thresholdStr != "" {
		// Parse the threshold
		threshold, err := strconv.ParseFloat(thresholdStr, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "Invalid threshold format"}`))
			return
		}

		// Update threshold
		if err := svc.UpdateThreshold(id, threshold); err != nil {
			logger.Println("Error updating threshold: ", err)
			http.Error(w, `{"error": "Failed to update threshold"}`, http.StatusInternalServerError)
			return
		}
	} else {
		// Update chosen location
		if err := svc.SetChosenLocation(id); err != nil {
			logger.Println("Error setting chosen location:", err)
			http.Error(w, `{"error": "Failed to set chosen location"}`, http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Location updated"}`))
}

func SetChosenLocationHandler(w http.ResponseWriter, r *http.Request, logger *log.Logger, svc service.LocationService) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error": "Invalid location ID"}`, http.StatusBadRequest)
		return
	}

	if err := svc.SetChosenLocation(id); err != nil {
		logger.Println("Error setting chosen location:", err)
		http.Error(w, `{"error": "Failed to set chosen location"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Chosen location updated"}`))
}

func DeleteHandler(w http.ResponseWriter, r *http.Request, logger *log.Logger, svc service.LocationService) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error": "Invalid location ID"}`, http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	aff, err := svc.DeleteLocation(&models.Location{ID: int64(id)}, ctx)
	if err != nil {
		logger.Println("Could not delete location:", err, id)
		http.Error(w, "Internal Server error", http.StatusInternalServerError)
		return
	}

	// * Check if the data was found and deleted
	if aff == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "Resource not found."}`))
		return
	}

	// * This is a Success, response in JSON and with a 204 status code when location was successfully deleted
	w.WriteHeader(http.StatusNoContent)
}
