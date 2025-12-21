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

// * The GET method retrieves a resource identified by a URI *
// * curl -X GET http://127.0.0.1:8080/data/1 -i -u admin:password -H "Content-Type: application/json"
func GetByIDHandler(w http.ResponseWriter, r *http.Request, logger *log.Logger, ds service.DataService) {
	rawID := r.PathValue("id")

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if numID, err := strconv.Atoi(rawID); err == nil {
		data, err := ds.ReadOne(numID, ctx)
		//id, err := strconv.Atoi(r.PathValue("id"))

		if err != nil {
			switch err.(type) {
			case service.DataError:
				// * If the error is a DataError, handle it as a client error
				w.WriteHeader(http.StatusBadRequest)
				//w.Write([]byte(`{"error": "` + err.Error() + `"}`))
				w.Write([]byte(`{"error": "Missconfigured ID."}`))
				return
			default:
				// * If it is not a DataError, handle it as a server error
				logger.Println("Error creating data:", err, data)
				http.Error(w, "Internal server error.", http.StatusInternalServerError)
				return
			}
		} else {
			if data == nil {
				// * This is a User Error, response in JSON and with a 404 status code
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(`{"error": "Resource not found."}`))
				return
			}
		}
		//logger.Println("Received GET /api/data/{id} from Arduino:")
		//logger.Printf("%+v\n", data)

		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(data); err != nil {
			logger.Println("Error encoding data:", err, data)
			http.Error(w, "Internal server error.", http.StatusInternalServerError)
			return
		}
	} else {
		data, err := ds.ReadLatest(rawID, ctx)
		if err != nil {
			// * This is a User Error: format of id is invalid, response in JSON and with a 400 status code
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "Missconfigured ID."}`))
			return
		}
		if data == nil {
			// * This is a User Error, response in JSON and with a 404 status code
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error": "Resource not found."}`))
			return
		}

		//logger.Println("Received GET /api/data/{id} from Arduino:")
		//logger.Printf("%+v\n", data)

		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(data); err != nil {
			logger.Println("Error encoding data:", err, data)
			http.Error(w, "Internal server error.", http.StatusInternalServerError)
			return
		}
	}
}
