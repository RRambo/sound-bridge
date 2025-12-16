package data_test

import (
	handlers "goapi/internal/api/handlers/data"
	service "goapi/internal/api/service/data"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDeleteInvalidID(t *testing.T) {
	mockDataService := &service.MockDataServiceSuccessful{}

	req, err := http.NewRequest("DELETE", "/data/invalid", nil)
	if err != nil {
		t.Fatal(err)
	}

	// If your router doesn't have SetPathValue, remove or adjust this line
	// req.SetPathValue("id", "invalid")

	rr := httptest.NewRecorder()
	handlers.DeleteHandler(rr, req, log.Default(), mockDataService)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	expected := `{"error": "Missconfigured ID."}`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestDeleteError(t *testing.T) {
	mockDataService := &service.MockDataServiceError{}

	req, err := http.NewRequest("DELETE", "/data/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.SetPathValue("id", "1")

	rr := httptest.NewRecorder()
	handlers.DeleteHandler(rr, req, log.Default(), mockDataService)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}

	expected := `Internal Server error`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestDeleteNotFound(t *testing.T) {
	mockDataService := &service.MockDataServiceNotFound{}

	req, err := http.NewRequest("DELETE", "/data/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.SetPathValue("id", "1")

	rr := httptest.NewRecorder()
	handlers.DeleteHandler(rr, req, log.Default(), mockDataService)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}

	expected := `{"error": "Resource not found."}`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestDeleteSuccessful(t *testing.T) {
	mockDataService := &service.MockDataServiceSuccessful{}

	req, err := http.NewRequest("DELETE", "/data/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.SetPathValue("id", "1")

	rr := httptest.NewRecorder()
	handlers.DeleteHandler(rr, req, log.Default(), mockDataService)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}

	if rr.Body.String() != "" {
		t.Errorf("handler returned unexpected body: got %v want empty body", rr.Body.String())
	}
}
