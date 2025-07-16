package service

import (
	"context"
	"fmt"
	"golang-bootcamp-July/internal/domain"
	"golang-bootcamp-July/internal/repository"
	"math/rand"
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
		return &domain.BookingResult{
			UserID:  userID,
			Success: false,
			Error:   domain.ErrInvalidUserID.Error(),
		}, domain.ErrInvalidUserID
	}

	time.Sleep(time.Microsecond * time.Duration(rand.Intn(100)))

	availableTickets, err := s.ticketRepo.GetAvailableTickets(ctx)
	if err != nil {
		return &domain.BookingResult{
			UserID:  userID,
			Success: false,
			Error:   err.Error(),
		}, err
	}

	if availableTickets <= 0 {
		return &domain.BookingResult{
			UserID:  userID,
			Success: false,
			Error:   domain.ErrNoTicketsAvailable.Error(),
		}, domain.ErrNoTicketsAvailable
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	availableTickets, err = s.ticketRepo.GetAvailableTickets(ctx)
	if err != nil {
		return &domain.BookingResult{
			UserID:  userID,
			Success: false,
			Error:   err.Error(),
		}, err
	}

	if availableTickets <= 0 {
		return &domain.BookingResult{
			UserID:  userID,
			Success: false,
			Error:   domain.ErrNoTicketsAvailable.Error(),
		}, domain.ErrNoTicketsAvailable
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
		return &domain.BookingResult{
			UserID:  userID,
			Success: false,
			Error:   err.Error(),
		}, err
	}

	// Decrement available tickets
	if err := s.ticketRepo.DecrementAvailableTickets(ctx, 1); err != nil {
		// Rollback booking if ticket decrement fails
		_ = s.bookingRepo.DeleteBooking(ctx, bookingID)
		return &domain.BookingResult{
			UserID:  userID,
			Success: false,
			Error:   err.Error(),
		}, err
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
		return &domain.BookingStats{}, err
	}

	availableTickets, err := s.ticketRepo.GetAvailableTickets(ctx)
	if err != nil {
		return &domain.BookingStats{}, err
	}

	bookedTickets := totalTicket - availableTickets

	return &domain.BookingStats{
		TotalTickets:     totalTicket,
		BookedTickets:    bookedTickets,
		AvailableTickets: availableTickets,
	}, nil
}

func (s *BookingService) GetUserBookingsService(ctx context.Context, userID string) ([]*domain.Booking, error) {
	return s.bookingRepo.GetBookingsByUserID(ctx, userID)
}
