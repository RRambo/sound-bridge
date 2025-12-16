package middleware

import (
	"encoding/base64"
	"net/http"
	"strings"
)

func BasicAuthenticationMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// * If type is Option return
		if r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
			return
		}

		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error": "Unauthorized: Missing credentials."}`))
			return
		}

		// * Split the Authorization header to get the 'Basic' part and the encoded credentials part
		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Basic" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "Malformed or invalid Authorization header. [1]"}`))
			return
		}

		// * Decode the credentials part of the header
		decoded, err := base64.StdEncoding.DecodeString(headerParts[1])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "Malformed or invalid Authorization header. [2]"}`))
			return
		}

		// * Split the decoded credentials to get the username and password
		credentials := strings.SplitN(string(decoded), ":", 2)
		if len(credentials) != 2 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "Malformed or invalid Authorization header. [3]"}`))
			return
		}

		username, password := credentials[0], credentials[1]

		if !validateUser(username, password) {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error": "Unauthorized: Invalid credentials."}`))
			return
		}

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}

func validateUser(username, password string) bool {
	/*username = os.Getenv("BASIC_USER")
	password = os.Getenv("BASIC_PASS")*/
	return username == "kids_noisemeter_admin" && password == "passwordkids"
}
