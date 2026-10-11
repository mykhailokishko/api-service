package domain

import "strings"

type PaymentSystem string

const (
	PaymentSystemUnknown    PaymentSystem = "Unknown"
	PaymentSystemUnion      PaymentSystem = "Union"
	PaymentSystemVisa       PaymentSystem = "Visa"
	PaymentSystemMastercard PaymentSystem = "Mastercard"
	PaymentSystemAmex       PaymentSystem = "American Express"
)

func DetectPaymentSystem(number string) PaymentSystem {
	n := len(number)
	if n < 13 || n > 19 {
		return PaymentSystemUnknown
	}

	if strings.HasPrefix(number, "4") && (n == 13 || n == 16 || n == 19) {
		return PaymentSystemVisa
	}

	if n == 15 && (strings.HasPrefix(number, "34") || strings.HasPrefix(number, "37")) {
		return PaymentSystemAmex
	}

	if n == 16 {
		prefixOfTwo := number[:2]
		if prefixOfTwo >= "51" && prefixOfTwo <= "55" {
			return PaymentSystemMastercard
		}
	}
	return PaymentSystemUnknown
}
