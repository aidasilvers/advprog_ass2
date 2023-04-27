package main

import (
	"encoding/json"
	"greenlight.bcc/internal/assert"
	"net/http/httptest"
	"testing"
)

type inputData struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Activated bool   `json:"activated"`
}

func Test_application_registerUserHandler(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routesTest())
	defer ts.Close()

	tests := []struct {
		name     string
		model    inputData
		wantCode int
	}{
		{
			name: "Test1",
			model: inputData{
				ID:        1,
				Name:      "test",
				Email:     "tes@mail.ru",
				Password:  "1234567",
				Activated: false,
			},
			wantCode: 404,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputData := struct {
				ID        int64  `json:"id"`
				Name      string `json:"name"`
				Email     string `json:"email"`
				Password  string `json:"password"`
				Activated bool   `json:"activated"`
			}{
				ID:        1,
				Name:      "Aida",
				Email:     "test@mail.ru",
				Password:  "123456789",
				Activated: true,
			}

			b, err := json.Marshal(&inputData)
			if err != nil {
				t.Fatal("wrong input data")
			}
			if tt.name == "test for wrong input" {
				b = append(b, 'a')
			}

			code, _, _ := ts.postForm(t, "/v1/users", b)

			assert.Equal(t, code, tt.wantCode)
		})
	}
}

func TestActivateUserHandler(t *testing.T) {
	app := newTestApplication(t)

	r := httptest.NewRequest("PUT", "/", nil)
	w := httptest.NewRecorder()
	app.activateUserHandler(w, r)
}
