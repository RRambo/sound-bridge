package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCommonInvalidContentType(t *testing.T) {

	// * Our id does not need to be valid, as we are only testing the content type *
	req, err := http.NewRequest("GET", "/data/0", nil)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "text/plain")
	rr := httptest.NewRecorder()

	handler := CommonMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnsupportedMediaType {
		t.Fatalf("Expected status code 415, got: %d", rr.Code)
	}

	expected := `{"error": "Content-Type header should be set to: application/json."}`
	if rr.Body.String() != expected {
		t.Fatalf("Expected response body: %s, got: %s", expected, rr.Body.String())
	}
}

func TestCommonCorrectContentType(t *testing.T) {

	req, err := http.NewRequest("GET", "/data/0", nil)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler := CommonMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	handler.ServeHTTP(rr, req)

	// * The status code is set in the handler, so we will not check it here; The handlers are tested separately *

	if rr.Header().Get("Content-Type") != "application/json" {
		t.Fatalf("Expected Content-Type: application/json, got: %s", rr.Header().Get("Content-Type"))
	}
	if rr.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Fatalf("Expected Access-Control-Allow-Origin: *, got: %s", rr.Header().Get("Access-Control-Allow-Origin"))
	}
}
