package domain

import "time"

type Booking struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Timestamp time.Time `json:"timestamp"`
}

type BookingResult struct {
	UserID    string `json:"user_id"`
	Success   bool   `json:"success"`
	BookingID string `json:"booking_id,omitempty"`
	Error     string `json:"error,omitempty"`
}

type BookingRequest struct {
	UserID string `json:"user_id"`
}

type BookingStats struct {
	TotalTickets     int32 `json:"total_tickets"`
	BookedTickets    int32 `json:"booked_tickets"`
	AvailableTickets int32 `json:"available_tickets"`
}
