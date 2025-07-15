package repository

import (
	"context"
	"golang-bootcamp-July/internal/domain"
)

type BookingRepository interface {
	CreateBooking(ctx context.Context, booking *domain.Booking) error

	GetBooking(ctx context.Context, id string) (*domain.Booking, error)

	GetBookingsByUserID(ctx context.Context, userID string) ([]*domain.Booking, error)

	GetAllBookings(ctx context.Context) ([]*domain.Booking, error)

	DeleteBooking(ctx context.Context, id string) error

	GetBookingCount(ctx context.Context) (int32, error)
}

type TicketRepository interface {
	GetAvailableTickets(ctx context.Context) (int32, error)

	DecrementAvailableTickets(ctx context.Context, count int32) error

	GetTotalTickets(ctx context.Context) (int32, error)

	SetTotalTickets(ctx context.Context, total int32) error
}
