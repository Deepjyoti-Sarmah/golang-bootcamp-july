package main

import (
	"context"
	"golang-bootcamp-July/internal/repository/inmemory"
	"golang-bootcamp-July/internal/service"
	"golang-bootcamp-July/internal/worker"
	"log"
	"strings"
	"sync/atomic"
	"time"
)

func main() {
	const (
		totalTickets  = 50_000
		totalRequests = 75_000
		numWorkers    = 1000
	)

	log.Println("🎫 Ticket Booking System Starting...")
	log.Printf("📊 Total Tickets Available: %d", totalTickets)
	log.Printf("👥 Total Booking Requests: %d", totalRequests)
	log.Printf("⚙️  Worker Pool Size: %d", numWorkers)
	log.Println(strings.Repeat("-", 50))

	// startTime := time.Now()

	// initialize repository
	bookingRepo := inmemory.NewBookingRepository()
	ticketRepo := inmemory.NewTicketRepository(totalTickets)

	bookingService := service.NewBookingService(bookingRepo, ticketRepo)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	workerPool := worker.NewPool(numWorkers, bookingService)
	workerPool.Start(ctx)

	var (
		successfulBookings int32
		failedBookings     int32
		processedRequests  int32
	)

	resultDone := make(chan bool)

	go func() {
		for result := range workerPool.Results() {
			atomic.AddInt32(&processedRequests, 1)

			if result.Success {
				atomic.AddInt32(&successfulBookings, 1)
				if processedRequests%5000 == 0 {
					log.Printf("✅ Progress: %d/%d requests processed, %d tickets booked",
						atomic.LoadInt32(&processedRequests),
						totalRequests,
						atomic.LoadInt32(&successfulBookings),
					)
				} else {
					atomic.AddInt32(&failedBookings, 1)
				}
			}
			resultDone <- true
		}
	}()

	log.Println("🚀 Starting concurrent booking requests...")

	// Client simulation
}
