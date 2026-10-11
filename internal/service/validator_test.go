package service_test

import (
	"errors"
	"testing"
	"time"

	"card-validator/internal/domain"
	"card-validator/internal/service"
)

func TestValidator_Validate(t *testing.T) {
	// Фіксуємо час: 15 жовтня 2026 року (детерміноване середовище)
	mockClock := domain.NewMockClock(time.Date(2026, 10, 15, 12, 0, 0, 0, time.UTC))
	validator := service.NewValidator(mockClock)

	tests := []struct {
		name                 string
		cardNumber           string
		month                int
		year                 int
		expectedErr          error
		expectedValid        bool
		expectedSystem       domain.PaymentSystem
	}{
		{
			name:           "Valid Visa card (Happy path)",
			cardNumber:     "4111111111111111",
			month:          12,
			year:           2028,
			expectedErr:    nil,
			expectedValid:  true,
			expectedSystem: domain.PaymentSystemVisa,
		},
		{
			name:           "Valid Mastercard (current month Oct 2026)",
			cardNumber:     "5555555555554444",
			month:          10,
			year:           2026,
			expectedErr:    nil,
			expectedValid:  true,
			expectedSystem: domain.PaymentSystemMastercard,
		},
		{
			name:           "Valid Amex card (15 digits)",
			cardNumber:     "378282246310005",
			month:          5,
			year:           2027,
			expectedErr:    nil,
			expectedValid:  true,
			expectedSystem: domain.PaymentSystemAmex,
		},
		{
			name:          "Format error: contains letters (005)",
			cardNumber:    "4111abcd11111111",
			month:         12,
			year:          2028,
			expectedErr:   domain.ErrInvalidCardFormat,
			expectedValid: false,
		},
		{
			name:          "Format error: too short (005)",
			cardNumber:    "411111111111",
			month:         12,
			year:          2028,
			expectedErr:   domain.ErrInvalidCardFormat,
			expectedValid: false,
		},
		{
			name:          "Luhn checksum error: from Solidgate task (002)",
			cardNumber:    "1111111111111111",
			month:         10,
			year:          2028,
			expectedErr:   domain.ErrInvalidCardNumber,
			expectedValid: false,
		},
		{
			name:          "Invalid month: month 13 (003)",
			cardNumber:    "4111111111111111",
			month:         13,
			year:          2028,
			expectedErr:   domain.ErrInvalidMonth,
			expectedValid: false,
		},
		{
			name:          "Invalid month: month 0 (003)",
			cardNumber:    "4111111111111111",
			month:         0,
			year:          2028,
			expectedErr:   domain.ErrInvalidMonth,
			expectedValid: false,
		},
		{
			name:          "Invalid year: year < 1000 (004)",
			cardNumber:    "4111111111111111",
			month:         12,
			year:          26,
			expectedErr:   domain.ErrInvalidYear,
			expectedValid: false,
		},
		{
			name:          "Expired card: from Solidgate task 01/2021 (001)",
			cardNumber:    "4111111111111111",
			month:         1,
			year:          2021,
			expectedErr:   domain.ErrCardExpired,
			expectedValid: false,
		},
		{
			name:          "Expired card: previous month of current year 09/2026 (001)",
			cardNumber:    "4111111111111111",
			month:         9,
			year:          2026,
			expectedErr:   domain.ErrCardExpired,
			expectedValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			card, err := validator.Validate(tt.cardNumber, tt.month, tt.year)

			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("Validate() error = %v, want %v", err, tt.expectedErr)
			}

			if tt.expectedValid {
				if card == nil {
					t.Fatalf("Validate() returned nil card, expected valid card")
				}
				if card.PaymentSystem != tt.expectedSystem {
					t.Errorf("Validate() PaymentSystem = %v, want %v", card.PaymentSystem, tt.expectedSystem)
				}
			} else {
				if card != nil {
					t.Errorf("Validate() returned card %v, expected nil", card)
				}
			}
		})
	}
}
