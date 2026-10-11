// Містить базові бізнес-сутності, правила та помилки платіжного домену.
package domain

import "fmt"

const (
	ErrCodeCardExpired       = "001"
	ErrCodeInvalidCardNumber = "002"
	ErrCodeInvalidMonth      = "003"
	ErrCodeInvalidYear       = "004"
	ErrCodeInvalidCardFormat = "005"
	ErrCodeMalformedPayload  = "006"
)

type DomainError struct {
	Code    string
	Message string
}

func (dm *DomainError) Error() string {
	return fmt.Sprintf("[%s] %s", dm.Code, dm.Message)

}

var (
	ErrCardExpired = &DomainError{
		Code:    ErrCodeCardExpired,
		Message: "Card has expired",
	}
	ErrInvalidCardNumber = &DomainError{
		Code:    ErrCodeInvalidCardNumber,
		Message: "Invalid card number",
	}
	ErrInvalidMonth = &DomainError{
		Code:    ErrCodeInvalidMonth,
		Message: "Invalid expiration month",
	}
	ErrInvalidYear = &DomainError{
		Code:    ErrCodeInvalidYear,
		Message: "Invalid expiration year",
	}
	ErrInvalidCardFormat = &DomainError{
		Code:    ErrCodeInvalidCardFormat,
		Message: "Invalid card format",
	}
	ErrMalformedPayload = &DomainError{
		Code:    ErrCodeMalformedPayload,
		Message: "Malformed payload",
	}
)
