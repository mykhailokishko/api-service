package domain

import (
	"errors"
	"testing"
)

func TestCard_ValidateFormat(t *testing.T) {
	tests := []struct {
		name        string
		number      string
		expectedErr error
	}{
		{name: "Valid 16 digits", number: "4111111111111111", expectedErr: nil},
		{name: "Valid 13 digits", number: "4111111111111", expectedErr: nil},
		{name: "Valid 19 digits", number: "6759123456789012342", expectedErr: nil},
		{name: "Too short (<13 digits)", number: "123456789012", expectedErr: ErrInvalidCardFormat},
		{name: "Too long (>19 digits)", number: "12345678901234567890", expectedErr: ErrInvalidCardFormat},
		{name: "Contains letters", number: "4111abcd11111111", expectedErr: ErrInvalidCardFormat},
		{name: "Contains dashes", number: "4111-1111-1111-1111", expectedErr: ErrInvalidCardFormat},
		{name: "Contains spaces", number: "4111 1111 1111 1111", expectedErr: ErrInvalidCardFormat},
		{name: "Empty string", number: "", expectedErr: ErrInvalidCardFormat},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			card := Card{Number: tt.number}
			err := card.ValidateFormat()
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("ValidateFormat() error = %v, want %v", err, tt.expectedErr)
			}
		})
	}
}
