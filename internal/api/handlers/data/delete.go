package data

import (
	"context"
	"goapi/internal/api/repository/models"
	service "goapi/internal/api/service/data"
	"log"
	"net/http"
	"strconv"
	"time"
)

// * The DELETE method removes a resource identified by a URI *
// * curl -X DELETE http://127.0.0.1:8080/data/1 -i -u admin:password -H "Content-Type: application/json"
func DeleteHandler(w http.ResponseWriter, r *http.Request, logger *log.Logger, ds service.DataService) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		// * This is a User Error: format of id is invalid, response in JSON and with a 400 status code
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Missconfigured ID."}`))
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	aff, err := ds.Delete(&models.Data{ID: id}, ctx)
	if err != nil {
		logger.Println("Could not delete data:", err, id)
		http.Error(w, "Internal Server error", http.StatusInternalServerError)
		return
	}

	// * Check if the data was found and deleted
	if aff == 0 {
		logger.Println("Could not delete data:", err, id)
		http.Error(w, `{"error": "Resource not found."}`, http.StatusNotFound)
		return
	}

	// * This is a Success, response in JSON and with a 204 status code when data was successfully deleted
	w.WriteHeader(http.StatusNoContent)
}
