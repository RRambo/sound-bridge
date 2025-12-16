package data

import (
	"context"
	"encoding/json"
	service "goapi/internal/api/service/data"
	"log"
	"net/http"
	"strconv"
	"time"
)

// * The GET method retrieves all resources identified by a URI *
// * curl -X GET http://127.0.0.1:8080/data -i -u admin:password -H "Content-Type: application/json"
func GetHandler(w http.ResponseWriter, r *http.Request, logger *log.Logger, ds service.DataService) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		if err == err.(*strconv.NumError) {
			// * There was no page specified, so we default to 0 *
			page = 0
		} else {
			// * Invalid page specified, return a 400 status code *
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "Invalid page specified."}`))
			return
		}
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	data, err := ds.ReadMany(page, 10, ctx)
	if err != nil {
		logger.Println("Could not get data:", err, data)
		http.Error(w, "Internal Server error.", http.StatusInternalServerError)
		return
	}
	if len(data) == 0 {
		logger.Println("Could not get data:", err, data)
		http.Error(w, `{"error": "Resource not found."}`, http.StatusNotFound)
		return
	}

	// * Return the data to the user as JSON with a 200 OK status code
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger.Println("Error encoding data:", err, data)
		http.Error(w, "Internal Server error.", http.StatusInternalServerError)
		return
	}
}
