package domain

import "errors"

var (
	ErrNoTicketsAvailable = errors.New("no tickets available")
	ErrorBookingNotFound  = errors.New("booking not found")
	ErrInvalidUserID      = errors.New("invalid user ID")
	ErrServerUnavailable  = errors.New("server unavilable")
)
