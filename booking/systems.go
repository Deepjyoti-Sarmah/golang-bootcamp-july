package booking

import (
	"fmt"
	"golang-bootcamp-July/model"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

type TicketSystem struct {
	totalTickets     int32
	availableTickets int32
	mu               sync.Mutex
	bookings         map[string]*model.Booking
}

func NewTicketSystem(totalTickets int32) *TicketSystem {
	return &TicketSystem{
		totalTickets:     totalTickets,
		availableTickets: totalTickets,
		bookings:         make(map[string]*model.Booking),
	}
}

func (ts *TicketSystem) BookTicket(userID string) model.BookingResult {
	// network latency
	time.Sleep(time.Microsecond * time.Duration(rand.Intn(100)))

	// non-blocking read
	if atomic.LoadInt32(&ts.availableTickets) <= 0 {
		return model.BookingResult{
			UserID:  userID,
			Success: false,
			Error:   fmt.Errorf("no tickets available"),
		}
	}

	// locking
	ts.mu.Lock()
	defer ts.mu.Unlock()

	// double check
	if ts.availableTickets <= 0 {
		return model.BookingResult{
			UserID:  userID,
			Success: false,
			Error:   fmt.Errorf("no tickets available"),
		}
	}

	// create
	bookingID := fmt.Sprintf("BOOK-%s-%d", userID, time.Now().UnixNano())
	booking := &model.Booking{
		ID:        bookingID,
		UserID:    userID,
		Timestamp: time.Now(),
	}

	// update
	ts.bookings[bookingID] = booking
	atomic.AddInt32(&ts.availableTickets, -1)

	return model.BookingResult{
		UserID:    userID,
		Success:   true,
		BookingID: bookingID,
	}
}

func (ts *TicketSystem) GetStats() (totalBooked int32, available int32) {
	available = atomic.LoadInt32(&ts.availableTickets)
	totalBooked = ts.totalTickets - available
	return
}
