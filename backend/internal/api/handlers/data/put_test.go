package data_test

import (
	"goapi/internal/api/handlers/data"
	service "goapi/internal/api/service/data"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPutInvalidRequestBody(t *testing.T) {
	mockDS := &service.MockDataServiceSuccessful{}

	req, err := http.NewRequest("PUT", "/data", io.NopCloser(strings.NewReader("Plain text, not JSON")))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	data.PutHandler(rr, req, log.Default(), mockDS)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("wrong status code: got %v, want %v", rr.Code, http.StatusBadRequest)
	}

	expected := `{"error": "Invalid request data. Please check your input."}`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("unexpected body: got %v, want %v", rr.Body.String(), expected)
	}
}

func TestPutHandlerError(t *testing.T) {
	mockDS := &service.MockDataServiceError{}

	jsonBody := `{"id":1,"device_id":"arduino_001","room_name":"PlayRoom_A","sound_level":75.5,"threshold":70.0,"measure_time":"2024-01-01T12:00:00Z","is_alert":true,"description":"Test update"}`
	req, err := http.NewRequest("PUT", "/data", io.NopCloser(strings.NewReader(jsonBody)))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	data.PutHandler(rr, req, log.Default(), mockDS)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("wrong status code: got %v, want %v", rr.Code, http.StatusInternalServerError)
	}

	expected := "Internal server error."
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("unexpected body: got %v, want %v", rr.Body.String(), expected)
	}
}

func TestPutDataNotFound(t *testing.T) {
	mockDS := &service.MockDataServiceNotFound{}

	jsonBody := `{"id":999,"device_id":"arduino_001","room_name":"PlayRoom_A","sound_level":75.5,"threshold":70.0,"measure_time":"2024-01-01T12:00:00Z","is_alert":false,"description":"Test update"}`
	req, err := http.NewRequest("PUT", "/data", io.NopCloser(strings.NewReader(jsonBody)))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	data.PutHandler(rr, req, log.Default(), mockDS)

	if rr.Code != http.StatusNotFound {
		t.Errorf("wrong status code: got %v, want %v", rr.Code, http.StatusNotFound)
	}

	expected := `{"error": "Resource not found."}`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("unexpected body: got %v, want %v", rr.Body.String(), expected)
	}
}

func TestPutHandlerSuccess(t *testing.T) {
	mockDS := &service.MockDataServiceSuccessful{}

	jsonBody := `{"id":1,"device_id":"arduino_001","room_name":"PlayRoom_B","sound_level":82.3,"threshold":70.0,"measure_time":"2024-10-27T14:30:00Z","is_alert":true,"description":"Success test"}`
	req, err := http.NewRequest("PUT", "/data", io.NopCloser(strings.NewReader(jsonBody)))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	data.PutHandler(rr, req, log.Default(), mockDS)

	if rr.Code != http.StatusOK {
		t.Errorf("wrong status code: got %v, want %v", rr.Code, http.StatusOK)
	}

	expected := `{"id":1,"device_id":"arduino_001","room_name":"PlayRoom_B","sound_level":82.3,"threshold":70,"measure_time":"2024-10-27T14:30:00Z","is_alert":true,"description":"Success test"}`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("unexpected body: got %v, want %v", rr.Body.String(), expected)
	}
}
