package domain

import "testing"

func TestDetectPaymentSystem(t *testing.T) {
	tests := []struct {
		name     string
		number   string
		expected PaymentSystem
	}{
		{name: "Visa 16 digits", number: "4111111111111111", expected: PaymentSystemVisa},
		{name: "Visa 13 digits", number: "4111111111111", expected: PaymentSystemVisa},
		{name: "Visa 19 digits", number: "4111111111111111111", expected: PaymentSystemVisa},
		{name: "Mastercard 51 prefix", number: "5111111111111111", expected: PaymentSystemMastercard},
		{name: "Mastercard 55 prefix", number: "5555555555554444", expected: PaymentSystemMastercard},
		{name: "Amex 34 prefix (15 digits)", number: "341111111111111", expected: PaymentSystemAmex},
		{name: "Amex 37 prefix (15 digits)", number: "378282246310005", expected: PaymentSystemAmex},
		{name: "Unknown brand", number: "1111111111111111", expected: PaymentSystemUnknown},
		{name: "Too short (<13 digits)", number: "4111", expected: PaymentSystemUnknown},
		{name: "Too long (>19 digits)", number: "41111111111111111111", expected: PaymentSystemUnknown},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DetectPaymentSystem(tt.number)
			if got != tt.expected {
				t.Errorf("DetectPaymentSystem(%s) = %s, want %s", tt.number, got, tt.expected)
			}
		})
	}
}
