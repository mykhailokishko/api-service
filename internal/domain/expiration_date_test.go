package domain

import (
	"errors"
	"testing"
	"time"
)

func TestExpirationDate_Validate(t *testing.T) {
	// Заморозка часу на 11 жовтня 2026 року:
	mockClock := NewMockClock(time.Date(2026, 10, 11, 12, 0, 0, 0, time.UTC))

	tests := []struct {
		name        string
		exp         ExpirationDate
		expectedErr error
	}{
		{
			name:        "Valid card in current month (October 2026)",
			exp:         ExpirationDate{Month: 10, Year: 2026},
			expectedErr: nil,
		},
		{
			name:        "Valid card in future (December 2028)",
			exp:         ExpirationDate{Month: 12, Year: 2028},
			expectedErr: nil,
		},
		{
			name:        "Expired card in past month of same year (September 2026)",
			exp:         ExpirationDate{Month: 9, Year: 2026},
			expectedErr: ErrCardExpired,
		},
		{
			name:        "Expired card in past year",
			exp:         ExpirationDate{Month: 1, Year: 2021},
			expectedErr: ErrCardExpired,
		},
		{
			name:        "Invalid month > 12",
			exp:         ExpirationDate{Month: 13, Year: 2026},
			expectedErr: ErrInvalidMonth,
		},
		{
			name:        "Invalid month < 1",
			exp:         ExpirationDate{Month: 0, Year: 2026},
			expectedErr: ErrInvalidMonth,
		},
		{
			name:        "Invalid year (< 1000)",
			exp:         ExpirationDate{Month: 10, Year: 999},
			expectedErr: ErrInvalidYear,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.exp.Validate(mockClock)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("Validate() error = %v, want %v", err, tt.expectedErr)
			}
		})
	}
}
