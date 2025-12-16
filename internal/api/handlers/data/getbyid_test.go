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

func TestGetByIDInvalidID(t *testing.T) {
	mockDS := &service.MockDataServiceSuccessful{}

	req, err := http.NewRequest("GET", "/data/invalid", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetPathValue("id", "invalid") // required for routing
	rr := httptest.NewRecorder()

	data.GetByIDHandler(rr, req, log.Default(), mockDS)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusBadRequest)
	}

	expected := `{"error": "Missconfigured ID."}`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestGetByIDInternalError(t *testing.T) {
	mockDS := &service.MockDataServiceError{}

	req, err := http.NewRequest("GET", "/data/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetPathValue("id", "1")
	rr := httptest.NewRecorder()

	data.GetByIDHandler(rr, req, log.Default(), mockDS)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusInternalServerError)
	}

	expected := `Internal server error.`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestGetByIDNotFound(t *testing.T) {
	mockDS := &service.MockDataServiceNotFound{}

	req, err := http.NewRequest("GET", "/data/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetPathValue("id", "1")
	rr := httptest.NewRecorder()

	data.GetByIDHandler(rr, req, log.Default(), mockDS)

	if rr.Code != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusNotFound)
	}

	expected := `{"error": "Resource not found."}`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestGetByIDSuccessful(t *testing.T) {
	mockDS := &service.MockDataServiceSuccessful{}

	req, err := http.NewRequest("GET", "/data/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetPathValue("id", "1")
	rr := httptest.NewRecorder()

	data.GetByIDHandler(rr, req, log.Default(), mockDS)

	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
	}

	// Avoid shadowing 'data' package by using 'mockData'
	mockData, _ := mockDS.ReadOne(1, req.Context())
	expected, _ := json.Marshal(mockData)

	if strings.TrimSpace(rr.Body.String()) != string(expected) {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), string(expected))
	}
}
