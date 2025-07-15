package memory

import (
	"context"
	"sync/atomic"
)

type TicketRepository struct {
	totalTickets     int32
	availableTickets int32
}

func NewTicketRepository(totalTickets int32) *TicketRepository {
	return &TicketRepository{
		totalTickets:     totalTickets,
		availableTickets: totalTickets,
	}
}

func (r *TicketRepository) GetAvailableTickets(ctx context.Context) (int32, error) {
	return atomic.LoadInt32(&r.availableTickets), nil
}

func (r *TicketRepository) DecrementAvailableTickets(ctx context.Context, count int32) error {
	atomic.AddInt32(&r.availableTickets, -count)
	return nil
}

func (r *TicketRepository) GetTotalTickets(ctx context.Context) (int32, error) {
	return atomic.LoadInt32(&r.totalTickets), nil
}

func (r *TicketRepository) SetTotalTickets(ctx context.Context, total int32) error {
	atomic.StoreInt32(&r.totalTickets, total)
	atomic.StoreInt32(&r.availableTickets, total)
	return nil
}
