package service

import (
	"card-validator/internal/domain"
	"card-validator/pkg/luhn"
)

type CardValidator interface {
	Validate(cardNumber string, month, year int) (*domain.Card, error)
}

type Validator struct {
	clock domain.Clock
}

func NewValidator(clock domain.Clock) *Validator {
	return &Validator{
		clock: clock,
	}
}

func (v *Validator) Validate(cardNumber string, month, year int) (*domain.Card, error) {
	card := domain.NewCard(cardNumber, domain.ExpirationDate{
		Month: month,
		Year:  year,
	})
	if err := card.ValidateFormat(); err != nil {
		return nil, err
	}
	if !luhn.Validate(cardNumber) {
		return nil, domain.ErrInvalidCardNumber
	}
	if err := card.ExpirationDate.Validate(v.clock); err != nil {
		return nil, err
	}
	return &card, nil
}
