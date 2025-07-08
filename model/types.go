package model

import (
	"time"
)

type Booking struct {
	ID        string
	UserID    string
	Timestamp time.Time
}

type BookingResult struct {
	UserID    string
	Success   bool
	BookingID string
	Error     error
}
