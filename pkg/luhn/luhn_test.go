package luhn

import "testing"

func TestValidate(t *testing.T) {
	tests := []struct {
		name     string
		number   string
		expected bool
	}{
		{
			name:     "Valid Visa card",
			number:   "4111111111111111",
			expected: true,
		},
		{
			name:     "Valid Amex card",
			number:   "378282246310005",
			expected: true,
		},
		{
			name:     "Valid Mastercard card",
			number:   "5555555555554444",
			expected: true,
		},
		{
			name:     "Valid Union/Maestro card",
			number:   "6759123456789012342",
			expected: true,
		},
		{
			name:     "Invalid card number (all '1')",
			number:   "1111111111111111",
			expected: false,
		},
		{
			name:     "Invalid card number (empty string)",
			number:   "",
			expected: false,
		},
		{
			name:     "Invalid card number (all '0')",
			number:   "0000000000000000",
			expected: false,
		},
		{
			name:     "Invalid card number (has character in it)",
			number:   "111111111a111111",
			expected: false,
		},
		{
			name:     "Invalid card number (with dashes)",
			number:   "1111-1111-1111-1111",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Validate(tt.number)
			if got != tt.expected {
				t.Errorf("Validate(%s) = %v; want %v", tt.number, got, tt.expected)
			}
		})
	}
}

func BenchmarkValidate(b *testing.B) {
	b.Run("Valid 16 digits (Happy path)", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			Validate("4111111111111111")
		}
	})

	b.Run("Valid 19 digits (Worse case)", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			Validate("6759123456789012342")
		}
	})

	b.Run("Invalid fail-fast", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			Validate("411111111111111A")
		}
	})
}
