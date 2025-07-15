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
	r.mu.RLock()
	defer r.mu.RUnlock()

	booking, exits := r.bookings[id]
	if !exits {
		return nil, domain.ErrorBookingNotFound
	}

	return booking, nil
}

func (r *BookingRepository) GetBookingsByUserID(ctx context.Context, userID string) ([]*domain.Booking, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var userBookings []*domain.Booking
	for _, booking := range r.bookings {
		if booking.UserID == userID {
			userBookings = append(userBookings, booking)
		}
	}

	return userBookings, nil
}

func (r *BookingRepository) GetAllBookings(ctx context.Context) ([]*domain.Booking, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	bookings := make([]*domain.Booking, 0, len(r.bookings))
	for _, booking := range r.bookings {
		bookings = append(bookings, booking)
	}

	return bookings, nil
}

func (r *BookingRepository) DeleteBooking(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.bookings[id]; !exists {
		return domain.ErrorBookingNotFound
	}

	delete(r.bookings, id)
	return nil
}

func (r *BookingRepository) GetBookingCount(ctx context.Context) (int32, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return int32(len(r.bookings)), nil
}
