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

// * When using PUT, the client sends a complete representation of a resource to replace the current version: Whole Resource Replacement. *
// * curl -X PUT http://127.0.0.1:8080/data -i -u admin:password -H "Content-Type: application/json" -d '{"id": 1, "content": "updated data"}'
func PutHandler(w http.ResponseWriter, r *http.Request, logger *log.Logger, ds service.DataService) {
	var data models.Data

	// * Decode the JSON payload from the request body into the data struct
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		// * This is a User Error: format of body is invalid, response in JSON and with a 400 status code
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Invalid request data. Please check your input."}`))
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	// * Try to update the data in the database
	if aff, err := ds.Update(&data, ctx); err != nil {
		switch err.(type) {
		case service.DataError:
			// * If the error is a DataError, handle it as a client error
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "` + err.Error() + `"}`))
			return
		default:
			// * If it is not a DataError, handle it as a server error
			logger.Println("Error updating data:", err, data)
			http.Error(w, "Internal server error.", http.StatusInternalServerError)
			return
		}
	} else if aff == 0 {
		// * This is a User Error, response in JSON and with a 404 status code
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "Resource not found."}`))
		return
	}

	// * Return the data to the user as JSON with a 200 OK status code
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger.Println("Error encoding data:", err, data)
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}
}
