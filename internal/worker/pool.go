package worker

import (
	"context"
	"golang-bootcamp-July/internal/domain"
	"golang-bootcamp-July/internal/service"
	"sync"
)

type Pool struct {
	workers        []*Worker
	requests       chan string
	results        chan domain.BookingResult
	wg             sync.WaitGroup
	bookingService *service.BookingService
}

func NewPool(numWorkers int, bookingService *service.BookingService) *Pool {
	pool := &Pool{
		requests:       make(chan string, numWorkers*2),
		results:        make(chan domain.BookingResult, numWorkers*10),
		workers:        make([]*Worker, numWorkers),
		bookingService: bookingService,
	}

	for i := 0; i < numWorkers; i++ {
		pool.workers[i] = NewWorker(i, pool.requests, pool.results, bookingService, &pool.wg)
	}

	return pool
}

func (p *Pool) Start(ctx context.Context) {
	for _, worker := range p.workers {
		p.wg.Add(1)
		go worker.Start(ctx)
	}
}

func (p *Pool) Stop() {
	close(p.requests)
	p.wg.Wait()
	close(p.results)
}

func (p *Pool) SubmitRequest(userID string) {
	p.requests <- userID
}

func (p *Pool) Results() <-chan domain.BookingResult {
	return p.results
}
