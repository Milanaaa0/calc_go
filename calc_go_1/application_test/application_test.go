package applicationtest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	application "github.com/pashapdev/calc_go/internal/application"
)

func TestCalculateHandler(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		body           application.Request
		expectedCode   int
		expectedResult *application.Response
		expectedError  *application.ErrorResponse
	}{
		{
			name:           "Valid Expression",
			method:         http.MethodPost,
			body:           application.Request{Expression: "1 + 1"},
			expectedCode:   http.StatusOK,
			expectedResult: &application.Response{Result: 2},
		},
		{
			name:           "Valid Expression with parentheses",
			method:         http.MethodPost,
			body:           application.Request{Expression: "(1 + 1) * 2"},
			expectedCode:   http.StatusOK,
			expectedResult: &application.Response{Result: 4},
		},
		{
			name:           "Valid Expression with multiple operations",
			method:         http.MethodPost,
			body:           application.Request{Expression: "1 + 2 * 3 - 4 / 2"},
			expectedCode:   http.StatusOK,
			expectedResult: &application.Response{Result: 5},
		},
		{
			name:         "Invalid JSON",
			method:       http.MethodPost,
			body:         application.Request{},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:          "Invalid Expression - invalid character",
			method:        http.MethodPost,
			body:          application.Request{Expression: "1 + a"},
			expectedCode:  http.StatusUnprocessableEntity,
			expectedError: &application.ErrorResponse{Error: "Expression is not valid"},
		},
		{
			name:          "Invalid Expression - unmatched parentheses",
			method:        http.MethodPost,
			body:          application.Request{Expression: "(1 + 2"},
			expectedCode:  http.StatusUnprocessableEntity,
			expectedError: &application.ErrorResponse{Error: "Expression is not valid"},
		},
		{
			name:          "Division by Zero",
			method:        http.MethodPost,
			body:          application.Request{Expression: "1 / 0"},
			expectedCode:  http.StatusInternalServerError,
			expectedError: &application.ErrorResponse{Error: "Internal server error"},
		},
		{
			name:          "Other Calculation Error", // Добавим тест для другого типа ошибки
			method:        http.MethodPost,
			body:          application.Request{Expression: "1 +"}, // Пример ошибки из calculate
			expectedCode:  http.StatusInternalServerError,
			expectedError: &application.ErrorResponse{Error: "Internal server error"},
		},

		{
			name:         "Method Not Allowed",
			method:       http.MethodGet,
			body:         application.Request{Expression: "1 + 1"},
			expectedCode: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body bytes.Buffer
			if tt.method == http.MethodPost {
				json.NewEncoder(&body).Encode(tt.body)
			}

			req := httptest.NewRequest(tt.method, "/calculate", &body)
			rr := httptest.NewRecorder()

			application.CalculateHandler(rr, req)

			if status := rr.Code; status != tt.expectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedCode)
			}

			if tt.expectedResult != nil {
				var result application.Response
				if err := json.Unmarshal(rr.Body.Bytes(), &result); err != nil {
					t.Errorf("failed to unmarshal response: %v", err)
				}
				if result != *tt.expectedResult {
					t.Errorf("handler returned unexpected result: got %+v want %+v", result, *tt.expectedResult)
				}
			}

			if tt.expectedError != nil {
				var errResp application.ErrorResponse
				if err := json.Unmarshal(rr.Body.Bytes(), &errResp); err != nil {
					t.Errorf("failed to unmarshal error response: %v", err)
				}
				if errResp != *tt.expectedError {
					t.Errorf("handler returned unexpected error: got %+v want %+v", errResp, *tt.expectedError)
				}
			}
		})
	}
}
