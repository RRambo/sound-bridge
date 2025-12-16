package data_test

import (
	"encoding/json"
	"goapi/internal/api/handlers/data"
	service "goapi/internal/api/service/data"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetHandlerSuccessful(t *testing.T) {
	mockDataService := &service.MockDataServiceSuccessful{}

	req, err := http.NewRequest("GET", "/data", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	// Call the handler
	data.GetHandler(rr, req, log.Default(), mockDataService)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Compare response body to what MockDataService returns
	mockData, _ := mockDataService.ReadMany(0, 10, req.Context())
	expected, _ := json.Marshal(mockData)

	if strings.TrimSpace(rr.Body.String()) != string(expected) {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), string(expected))
	}
}

func TestGetHandlerNotFound(t *testing.T) {
	mockDataService := &service.MockDataServiceNotFound{}

	req, err := http.NewRequest("GET", "/data", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	data.GetHandler(rr, req, log.Default(), mockDataService)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}

	expected := `{"error": "Resource not found."}`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestGetHandlerError(t *testing.T) {
	mockDataService := &service.MockDataServiceError{}

	req, err := http.NewRequest("GET", "/data", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	data.GetHandler(rr, req, log.Default(), mockDataService)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}

	expected := `Internal Server error.`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
