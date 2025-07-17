package service

import (
	"context"
	"fmt"
	"golang-bootcamp-July/internal/domain"
	"golang-bootcamp-July/internal/repository"
	"sync"
	"time"
)

type BookingService struct {
	bookingRepo repository.BookingRepository
	ticketRepo  repository.TicketRepository
	mu          sync.Mutex
}

func NewBookingService(bookingRepo repository.BookingRepository, ticketRepo repository.TicketRepository) *BookingService {
	return &BookingService{
		bookingRepo: bookingRepo,
		ticketRepo:  ticketRepo,
	}
}

func (s *BookingService) BookTicketService(ctx context.Context, userID string) (*domain.BookingResult, error) {
	if userID == "" {
		return nil, domain.ErrInvalidUserID
	}

	// time.Sleep(time.Microsecond * time.Duration(rand.Intn(100)))

	availableTickets, err := s.ticketRepo.GetAvailableTickets(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get available tickets: %w ", err)
	}

	if availableTickets <= 0 {
		return nil, domain.ErrNoTicketsAvailable
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	availableTickets, err = s.ticketRepo.GetAvailableTickets(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get available tickets: %w ", err)
	}

	if availableTickets <= 0 {
		return nil, domain.ErrNoTicketsAvailable
	}

	// Create booking
	bookingID := fmt.Sprintf("BOOK-%s-%d", userID, time.Now().UnixNano())
	booking := &domain.Booking{
		ID:        bookingID,
		UserID:    userID,
		Timestamp: time.Now(),
	}

	// Save booking
	if err := s.bookingRepo.CreateBooking(ctx, booking); err != nil {
		return nil, fmt.Errorf("failed to create boooking: %w", err)
	}

	// Decrement available tickets
	if err := s.ticketRepo.DecrementAvailableTickets(ctx, 1); err != nil {
		// Rollback
		_ = s.bookingRepo.DeleteBooking(ctx, bookingID)
		return nil, fmt.Errorf("failed to decrement tickets: %w", err)
	}

	return &domain.BookingResult{
		UserID:    userID,
		Success:   true,
		BookingID: bookingID,
	}, nil
}

func (s *BookingService) GetBookingStatsService(ctx context.Context) (*domain.BookingStats, error) {
	totalTicket, err := s.ticketRepo.GetTotalTickets(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get total tickets: %w", err)
	}

	availableTickets, err := s.ticketRepo.GetAvailableTickets(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get available tickets: %w", err)
	}

	bookedTickets := totalTicket - availableTickets

	return &domain.BookingStats{
		TotalTickets:     totalTicket,
		BookedTickets:    bookedTickets,
		AvailableTickets: availableTickets,
	}, nil
}

func (s *BookingService) GetUserBookingsService(ctx context.Context, userID string) ([]*domain.Booking, error) {
	if userID == "" {
		return nil, domain.ErrInvalidUserID
	}

	bookings, err := s.bookingRepo.GetBookingsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user bookings: %w", err)
	}

	return bookings, nil
}
