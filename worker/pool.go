package worker

import (
	"context"
	"golang-bootcamp-July/booking"
	"golang-bootcamp-July/model"
	"sync"
)

// worker to process request from a channel
type Worker struct {
	id       int
	requests <-chan string
	results  chan<- model.BookingResult
	ts       *booking.TicketSystem
	wg       *sync.WaitGroup
}

// start booking requests
func (w *Worker) StartWorker(ctx context.Context) {
	defer w.wg.Done()

	for {
		select {
		case userID, ok := <-w.requests:
			if !ok {
				// channel closed, worker should exit
				return
			}

			result := w.ts.BookTicket(userID)

			// send result
			select {
			case w.results <- result:
			case <-ctx.Done():
				return
			}

		case <-ctx.Done():
			return
		}
	}
}

type WorkerPool struct {
	workers  []*Worker
	requests chan string
	results  chan model.BookingResult
	wg       sync.WaitGroup
}

func NewWorkerPool(numWorkers int, ts *booking.TicketSystem) *WorkerPool {
	pool := &WorkerPool{
		requests: make(chan string, numWorkers*2),
		results:  make(chan model.BookingResult, numWorkers*10),
		workers:  make([]*Worker, numWorkers),
	}

	for i := 0; i < numWorkers; i++ {
		pool.workers[i] = &Worker{
			id:       i,
			requests: pool.requests,
			results:  pool.results,
			ts:       ts,
			wg:       &pool.wg,
		}
	}

	return pool
}

// start all the workers
func (p *WorkerPool) StartWorkerPool(ctx context.Context) {
	for _, worker := range p.workers {
		p.wg.Add(1)
		go worker.StartWorker(ctx)
	}
}

func (p *WorkerPool) Stop() {
	close(p.requests)
	p.wg.Wait()
	close(p.results)
}

func (p *WorkerPool) SubmitRequest(userID string) {
	p.requests <- userID
}

func (p *WorkerPool) Results() <-chan model.BookingResult {
	return p.results
}
