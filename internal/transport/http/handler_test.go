package http_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"card-validator/internal/domain"
	"card-validator/internal/service"
	transporthttp "card-validator/internal/transport/http"
)

func TestCardHandler_ValidateCard(t *testing.T) {
	// Зафіксований час: 11 жовтня 2026 року
	mockClock := domain.NewMockClock(time.Date(2026, 10, 11, 12, 0, 0, 0, time.UTC))
	validator := service.NewValidator(mockClock)
	handler := transporthttp.NewCardHandler(validator)

	tests := []struct {
		name           string
		method         string
		payload        string
		expectedStatus int
		expectedValid  bool
		expectedCode   string
	}{
		{
			name:           "Valid Visa card (Happy path)",
			method:         http.MethodPost,
			payload:        `{"card_number":"4111111111111111","expiration_month":12,"expiration_year":2028}`,
			expectedStatus: http.StatusOK,
			expectedValid:  true,
		},
		{
			name:           "Valid card with flexible types (number as int, dates as strings)",
			method:         http.MethodPost,
			payload:        `{"card_number":4111111111111111,"expiration_month":"10","expiration_year":"2026"}`,
			expectedStatus: http.StatusOK,
			expectedValid:  true,
		},
		{
			name:           "Invalid Luhn card from Solidgate task (all 1s)",
			method:         http.MethodPost,
			payload:        `{"card_number":"1111111111111111","expiration_month":10,"expiration_year":2028}`,
			expectedStatus: http.StatusOK,
			expectedValid:  false,
			expectedCode:   "002",
		},
		{
			name:           "Expired card from Solidgate task (01/2021)",
			method:         http.MethodPost,
			payload:        `{"card_number":"4111111111111111","expiration_month":1,"expiration_year":2021}`,
			expectedStatus: http.StatusOK,
			expectedValid:  false,
			expectedCode:   "001",
		},
		{
			name:           "Invalid expiration month (13)",
			method:         http.MethodPost,
			payload:        `{"card_number":"4111111111111111","expiration_month":13,"expiration_year":2028}`,
			expectedStatus: http.StatusOK,
			expectedValid:  false,
			expectedCode:   "003",
		},
		{
			name:           "Invalid format: letters in card number",
			method:         http.MethodPost,
			payload:        `{"card_number":"4111-abcd-1111","expiration_month":10,"expiration_year":2028}`,
			expectedStatus: http.StatusOK,
			expectedValid:  false,
			expectedCode:   "005",
		},
		{
			name:           "Malformed JSON payload",
			method:         http.MethodPost,
			payload:        `{"card_number": 4111, invalid json`,
			expectedStatus: http.StatusOK,
			expectedValid:  false,
			expectedCode:   "006",
		},
		{
			name:           "Missing required field",
			method:         http.MethodPost,
			payload:        `{"card_number":"4111111111111111"}`,
			expectedStatus: http.StatusOK,
			expectedValid:  false,
			expectedCode:   "006",
		},
		{
			name:           "Method Not Allowed: GET request",
			method:         http.MethodGet,
			payload:        "",
			expectedStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/validate", bytes.NewBufferString(tt.payload))
			rec := httptest.NewRecorder()

			handler.ValidateCard(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Fatalf("status code = %d, want %d", rec.Code, tt.expectedStatus)
			}

			if tt.expectedStatus == http.StatusOK {
				var resp transporthttp.ValidateCardResponse
				if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response JSON: %v", err)
				}

				if resp.Valid != tt.expectedValid {
					t.Errorf("resp.Valid = %v, want %v", resp.Valid, tt.expectedValid)
				}

				if !tt.expectedValid {
					if resp.Error == nil {
						t.Fatalf("resp.Error is nil, want error with code %s", tt.expectedCode)
					}
					if resp.Error.Code != tt.expectedCode {
						t.Errorf("error.code = %s, want %s", resp.Error.Code, tt.expectedCode)
					}
				}
			}
		})
	}
}

func TestCardHandler_RegisterRoutes(t *testing.T) {
	mockClock := domain.NewMockClock(time.Now())
	validator := service.NewValidator(mockClock)
	handler := transporthttp.NewCardHandler(validator)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	// Перевіряється ендпоінт /validate
	req := httptest.NewRequest(http.MethodPost, "/validate", bytes.NewBufferString(`{"card_number":"4111111111111111","expiration_month":12,"expiration_year":2028}`))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("/validate endpoint returned status %d, want 200", rec.Code)
	}

	// Альтернативний ендпоінт /api/v1/cards/validate
	reqV1 := httptest.NewRequest(http.MethodPost, "/api/v1/cards/validate", bytes.NewBufferString(`{"card_number":"4111111111111111","expiration_month":12,"expiration_year":2028}`))
	recV1 := httptest.NewRecorder()
	mux.ServeHTTP(recV1, reqV1)

	if recV1.Code != http.StatusOK {
		t.Errorf("/api/v1/cards/validate endpoint returned status %d, want 200", recV1.Code)
	}
}
