package http

import (
	"card-validator/internal/domain"
	"card-validator/internal/service"
	"encoding/json"
	"errors"
	"net/http"
)

type CardHandler struct {
	validator service.CardValidator
}

func NewCardHandler(validator service.CardValidator) *CardHandler {
	return &CardHandler{
		validator: validator,
	}
}

func (h *CardHandler) ValidateCard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	var req ValidateCardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(NewErrorResponse(domain.ErrCodeMalformedPayload, "Invalid request payload"))
		return
	}
	_, err := h.validator.Validate(req.CardNumber, req.ExpirationMonth, req.ExpirationYear)
	if err != nil {
		var domErr *domain.DomainError
		if errors.As(err, &domErr) {
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(NewErrorResponse(domErr.Code, domErr.Message))
			return
		}
		// Якщо сталася якась інша помилка
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(NewErrorResponse(domain.ErrCodeMalformedPayload, err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(NewValidResponse())
}

func (h *CardHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/validate", h.ValidateCard)
	mux.HandleFunc("/api/v1/cards/validate", h.ValidateCard)
}
