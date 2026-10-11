package http

// ValidateCardResponse представляє JSON-контракт відповіді згідно з ТЗ Solidgate.
type ValidateCardResponse struct {
	Valid bool         `json:"valid"`
	Error *ErrorDetail `json:"error,omitempty"`
}

// ErrorDetail описує помилку валідації.
type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// NewValidResponse створює успішну відповідь.
func NewValidResponse() ValidateCardResponse {
	return ValidateCardResponse{
		Valid: true,
	}
}

// NewErrorResponse створює відповідь із помилкою.
func NewErrorResponse(code, message string) ValidateCardResponse {
	return ValidateCardResponse{
		Valid: false,
		Error: &ErrorDetail{
			Code:    code,
			Message: message,
		},
	}
}
