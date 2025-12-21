package data_test

import (
	"goapi/internal/api/handlers/data"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOptions(t *testing.T) {

	req, err := http.NewRequest("OPTIONS", "/data", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	data.OptionsHandler(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if rr.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Errorf("handler returned unexpected header: got %v want %v", rr.Header().Get("Access-Control-Allow-Origin"), "*")
	}

	if rr.Header().Get("Access-Control-Allow-Methods") != "GET, POST, PUT, DELETE" {
		t.Errorf("handler returned unexpected header: got %v want %v", rr.Header().Get("Access-Control-Allow-Methods"), "GET, POST, PUT, DELETE")
	}

	if rr.Header().Get("Access-Control-Allow-Headers") != "Content-Type, Authorization" {
		t.Errorf("handler returned unexpected header: got %v want %v", rr.Header().Get("Access-Control-Allow-Headers"), "Content-Type, Authorization")
	}

	if rr.Body.String() != "" {
		t.Errorf("handler returned unexpected body: got %v want empty body", rr.Body.String())
	}
}
