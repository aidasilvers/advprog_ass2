package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthcheckHandler(t *testing.T) {
	app := &application{
		config: &config{env: "test"},
	}

	req, err := http.NewRequest("GET", "/healthcheck", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.healthcheckHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var resp envelope
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		t.Errorf("unable to parse response body: %v", err)
	}

	if resp["status"] != "available" {
		t.Errorf("unexpected status: got %v want %v", resp["status"], "available")
	}

	if resp["system_info"].(map[string]interface{})["environment"] != "test" {
		t.Errorf("unexpected environment: got %v want %v", resp["system_info"].(map[string]interface{})["environment"], "test")
	}

	if resp["system_info"].(map[string]interface{})["version"] != version {
		t.Errorf("unexpected version: got %v want %v", resp["system_info"].(map[string]interface{})["version"], version)
	}

	// You can add more checks for the response fields here, if needed.
}
