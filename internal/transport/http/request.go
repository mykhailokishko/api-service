package http

import (
	"bytes"
	"encoding/json"
	"strconv"
	"strings"

	"card-validator/internal/domain"
)

type ValidateCardRequest struct {
	CardNumber      string
	ExpirationMonth int
	ExpirationYear  int
}

func (r *ValidateCardRequest) UnmarshalJSON(data []byte) error {
	var raw struct {
		CardNumber      json.RawMessage `json:"card_number"`
		ExpirationMonth json.RawMessage `json:"expiration_month"`
		ExpirationYear  json.RawMessage `json:"expiration_year"`
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return domain.ErrMalformedPayload
	}

	if len(raw.CardNumber) == 0 || len(raw.ExpirationMonth) == 0 || len(raw.ExpirationYear) == 0 {
		return domain.ErrMalformedPayload
	}

	cardNum, err := parseStringOrNumber(raw.CardNumber)
	if err != nil {
		return domain.ErrMalformedPayload
	}
	r.CardNumber = cardNum

	month, err := parseIntOrString(raw.ExpirationMonth)
	if err != nil {
		return domain.ErrMalformedPayload
	}
	r.ExpirationMonth = month

	year, err := parseIntOrString(raw.ExpirationYear)
	if err != nil {
		return domain.ErrMalformedPayload
	}
	r.ExpirationYear = year

	return nil
}

func parseStringOrNumber(raw []byte) (string, error) {
	raw = bytes.TrimSpace(raw)

	if len(raw) == 0 {
		return "", domain.ErrMalformedPayload
	}

	if raw[0] == '"' {
		var s string
		if err := json.Unmarshal(raw, &s); err != nil {
			return "", err
		}

		return strings.TrimSpace(s), nil
	}
	return string(raw), nil
}

func parseIntOrString(raw []byte) (int, error) {
	raw = bytes.TrimSpace(raw)
	if len(raw) == 0 {
		return 0, domain.ErrMalformedPayload
	}
	if raw[0] == '"' {
		var s string
		if err := json.Unmarshal(raw, &s); err != nil {
			return 0, err
		}
		return strconv.Atoi(strings.TrimSpace(s))
	}
	return strconv.Atoi(string(raw))
}
