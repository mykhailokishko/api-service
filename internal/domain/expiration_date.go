package domain

import "time"

type ExpirationDate struct {
	Month int
	Year  int
}

func (ed *ExpirationDate) Validate(clock Clock) error {
	if ed.Month < 1 || ed.Month > 12 {
		return ErrInvalidMonth
	}

	now := clock.Now().UTC()

	if ed.Year < 1000 || ed.Year > now.Year()+20 {
		return ErrInvalidYear // 004
	}

	endOfMonth := time.Date(ed.Year, time.Month(ed.Month)+1, 1, 0, 0, 0, 0, time.UTC)
	if !now.Before(endOfMonth) {
		return ErrCardExpired
	}

	return nil
}
