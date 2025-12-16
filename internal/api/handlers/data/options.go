package data

import "net/http"

// * The OPTIONS method is used to describe the communication options for the target resource. *
// * curl -X OPTIONS http://127.0.0.1:8080/data -i
func OptionsHandler(w http.ResponseWriter, r *http.Request) {
	// Preflight request: server returns a 200 OK status code and the allowed methods and headers in the response headers.
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.WriteHeader(http.StatusOK)
}
