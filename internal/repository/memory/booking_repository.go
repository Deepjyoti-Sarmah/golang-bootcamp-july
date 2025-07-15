package memory

import (
	"context"
	"golang-bootcamp-July/internal/domain"
	"sync"
)

type BookingRepository struct {
	bookings map[string]*domain.Booking
	mu       sync.RWMutex
}

func NewBookingRepository() *BookingRepository {
	return &BookingRepository{
		bookings: make(map[string]*domain.Booking),
	}
}

func (r *BookingRepository) CreateBooking(ctx context.Context, booking *domain.Booking) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.bookings[booking.ID] = booking
	return nil
}

func (r *BookingRepository) GetBooking(ctx context.Context, id string) (*domain.Booking, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	booking, exits := r.bookings[id]
	if !exits {
		return nil, domain.ErrorBookingNotFound
	}

	return booking, nil
}
