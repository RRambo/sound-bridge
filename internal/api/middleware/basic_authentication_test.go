package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// * Test: No Authorization header
func TestBasicAuthMissingAuthHeader(t *testing.T) {

	req := httptest.NewRequest(http.MethodGet, "/data/0", nil)
	rr := httptest.NewRecorder()

	handler := BasicAuthenticationMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// * This should not be called if the Authorization header is missing,
		// * The Authorization header is checked in the middleware before calling the handler and if it is missing the handler should not be called
		t.Error("Handler should not have been called")
	}),
	)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, rr.Code)
	}

	// * Test if returned error message is correct
	expected := `{"error": "Unauthorized: Missing credentials."}`
	if rr.Body.String() != expected {
		t.Errorf("Expected body %s, got %s", expected, rr.Body.String())
	}
}

// * Test: Invalid Authorization header
func TestBasicAuthMalformedHeader(t *testing.T) {

	req := httptest.NewRequest(http.MethodGet, "/data/0", nil)
	req.Header.Add("Authorization", "INVALID")

	rr := httptest.NewRecorder()
	handler := BasicAuthenticationMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// * This should not be called if the Authorization header is missing,
		// * The Authorization header is checked in the middleware before calling the handler and if it is missing the handler should not be called
		t.Error("Handler should not have been called")
	}),
	)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rr.Code)
	}

	// * Test if returned error message is correct
	expected := `{"error": "Malformed or invalid Authorization header. [1]"}`
	if rr.Body.String() != expected {
		t.Errorf("Expected body %s, got %s", expected, rr.Body.String())
	}
}

func TestBasicAuthErrorOnBase64Decoding(t *testing.T) {

	req := httptest.NewRequest(http.MethodGet, "/data/0", nil)
	req.Header.Add("Authorization", "Basic XXX")

	rr := httptest.NewRecorder()
	handler := BasicAuthenticationMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// * This should not be called if the Authorization header is missing,
		// * The Authorization header is checked in the middleware before calling the handler and if it is missing the handler should not be called
		t.Error("Handler should not have been called")
	}),
	)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rr.Code)
	}

	// * Test if returned error message is correct
	expected := `{"error": "Malformed or invalid Authorization header. [2]"}`
	if rr.Body.String() != expected {
		t.Errorf("Expected body %s, got %s", expected, rr.Body.String())
	}
}

func TestBasicAuthErrorOnSplitDecodedString(t *testing.T) {

	req := httptest.NewRequest(http.MethodGet, "/data/0", nil)
	req.Header.Add("Authorization", "Basic RWluYXJUZXN0YWE=")

	rr := httptest.NewRecorder()
	handler := BasicAuthenticationMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// * This should not be called if the Authorization header is missing,
		// * The Authorization header is checked in the middleware before calling the handler and if it is missing the handler should not be called
		t.Error("Handler should not have been called")
	}),
	)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rr.Code)
	}

	// * Test if returned error message is correct
	expected := `{"error": "Malformed or invalid Authorization header. [3]"}`
	if rr.Body.String() != expected {
		t.Errorf("Expected body %s, got %s", expected, rr.Body.String())
	}
}

func TestBasicAuthInvalidCredentials(t *testing.T) {

	req := httptest.NewRequest(http.MethodGet, "/data/0", nil)
	req.Header.Add("Authorization", "Basic RWluYXI6RWluYXI=")

	rr := httptest.NewRecorder()
	handler := BasicAuthenticationMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// * This should not be called if the Authorization header is missing,
		// * The Authorization header is checked in the middleware before calling the handler and if it is missing the handler should not be called
		t.Error("Handler should not have been called")
	}),
	)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, rr.Code)
	}

	// * Test if returned error message is correct
	expected := `{"error": "Unauthorized: Invalid credentials."}`
	if rr.Body.String() != expected {
		t.Errorf("Expected body %s, got %s", expected, rr.Body.String())
	}

}
