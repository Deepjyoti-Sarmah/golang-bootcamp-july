package worker

import (
	"context"
	"golang-bootcamp-July/internal/domain"
	"golang-bootcamp-July/internal/service"
	"sync"
)

type Worker struct {
	id              int
	requests        <-chan string
	results         chan<- domain.BookingResult
	boookingService *service.BookingService
	wg              *sync.WaitGroup
}

func NewWorker(id int, requests <-chan string, results chan<- domain.BookingResult, bookingService *service.BookingService, wg *sync.WaitGroup) *Worker {
	return &Worker{
		id:              id,
		requests:        requests,
		results:         results,
		boookingService: bookingService,
		wg:              wg,
	}
}

func (w *Worker) Start(ctx context.Context) {
	defer w.wg.Done()

	for {
		select {
		case userID, ok := <-w.requests:
			if !ok {
				return
			}

			result, err := w.boookingService.BookTicketService(ctx, userID)
			if err != nil {
				select {
				case w.results <- domain.BookingResult{
					UserID:  userID,
					Success: false,
					Error:   err.Error(),
				}:
				case <-ctx.Done():
					return
				}
				continue
			}

			select {
			case w.results <- *result:
			case <-ctx.Done():
				return
			}

		case <-ctx.Done():
			return
		}
	}
}
