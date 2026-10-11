package domain

type Card struct {
	Number         string
	ExpirationDate ExpirationDate
	PaymentSystem  PaymentSystem
}

func NewCard(number string, exp ExpirationDate) Card {
	return Card{
		Number:         number,
		ExpirationDate: exp,
		PaymentSystem:  DetectPaymentSystem(number),
	}
}

func (c Card) ValidateFormat() error {
	n := len(c.Number)
	if n < 13 || n > 19 {
		return ErrInvalidCardFormat
	}
	for i := 0; i < n; i++ {
		if c.Number[i] < '0' || c.Number[i] > '9' {
			return ErrInvalidCardFormat
		}
	}
	return nil
}
