package main

import (
	"bytes"
	"encoding/json"
	"greenlight.bcc/internal/data"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateAuthToken(t *testing.T) {
	app := newTestApplication(t)
	user := &data.User{
		Email: "test@example.com",
	}
	err := app.models.Users.Insert(user)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		Name   string
		Email  string
		Status int
	}{
		{
			Name:   "Valid credentials",
			Email:  "test@example.com",
			Status: http.StatusCreated,
		},
		{
			Name:  "Invalid email",
			Email: "invalid",

			Status: http.StatusBadRequest,
		},
		{
			Name:  "Invalid password",
			Email: "test@example.com",

			Status: http.StatusBadRequest,
		},
		{
			Name:  "Incorrect email",
			Email: "incorrect@example.com",

			Status: http.StatusUnauthorized,
		},
		{
			Name:  "Incorrect password",
			Email: "test@example.com",

			Status: http.StatusUnauthorized,
		},
	}

	// Perform tests
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			body := map[string]string{
				"email": test.Email,
			}
			jsonBody, err := json.Marshal(body)
			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequest("POST", "/auth/token", bytes.NewBuffer(jsonBody))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()

			app.createAuthenticationTokenHandler(recorder, req)

			if recorder.Code != test.Status {
				t.Errorf("expected status %d but got %d", test.Status, recorder.Code)
			}

			if test.Status == http.StatusCreated {
				var responseBody map[string]string
				err = json.NewDecoder(recorder.Body).Decode(&responseBody)
				if err != nil {
					t.Fatal(err)
				}
				if _, ok := responseBody["authentication_token"]; !ok {
					t.Error("expected authentication_token in response body")
				}
			}
		})
	}
}
