package main

import (
	"context"
	"fmt"
	"golang-bootcamp-July/internal/repository/inmemory"
	"golang-bootcamp-July/internal/service"
	"golang-bootcamp-July/internal/worker"
	"log"
	"strings"
	"sync"
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

	startTime := time.Now()

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

	var resultWg sync.WaitGroup
	resultWg.Add(1)

	go func() {
		defer resultWg.Done()
		for result := range workerPool.Results() {
			atomic.AddInt32(&processedRequests, 1)
			processed := atomic.LoadInt32(&processedRequests)

			if result.Success {
				atomic.AddInt32(&successfulBookings, 1)
			} else {
				atomic.AddInt32(&failedBookings, 1)
			}

			if processed%5000 == 0 {
				log.Printf("✅ Progress: %d/%d requests processed, %d successful, %d failed",
					processed,
					totalRequests,
					atomic.LoadInt32(&successfulBookings),
					atomic.LoadInt32(&failedBookings),
				)
			}
		}
	}()

	log.Println("🚀 Starting concurrent booking requests...")

	// Client simulation
	go func() {
		var requestWg sync.WaitGroup

		batchSize := 1000
		for i := 0; i < totalRequests; i += batchSize {
			requestWg.Add(1)
			go func(start int) {
				defer requestWg.Done()
				end := start + batchSize
				if end > totalRequests {
					end = totalRequests
				}

				for j := start; j < end; j++ {
					userID := fmt.Sprintf("USER-%06d", j)
					workerPool.SubmitRequest(userID)
				}
			}(i)
		}

		requestWg.Wait()
		workerPool.Stop()
	}()

	resultWg.Wait()

	duration := time.Since(startTime)
	stats, _ := bookingService.GetBookingStatsService(ctx)

	log.Println(strings.Repeat("=", 50))
	log.Println("📈 BOOKING SYSTEM FINAL REPORT")
	log.Println(strings.Repeat("=", 50))
	log.Printf("⏱️  Total Time Taken: %v", duration)
	log.Printf("🎫 Total Tickets: %d", stats.TotalTickets)
	log.Printf("✅ Total Tickets Booked: %d", stats.BookedTickets)
	log.Printf("❌ Total Tickets NOT Booked: %d", atomic.LoadInt32(&failedBookings))
	log.Printf("📊 Total Requests Processed: %d", atomic.LoadInt32(&processedRequests))
	log.Printf("🎯 Remaining Tickets: %d", stats.AvailableTickets)
	log.Printf("⚡ Requests per Second: %.2f", float64(totalRequests)/duration.Seconds())
	log.Printf("🔄 Average Time per Request: %v", duration/time.Duration(totalRequests))
	log.Println(strings.Repeat("=", 50))

	if stats.BookedTickets == stats.TotalTickets && stats.AvailableTickets == 0 {
		log.Println("✅ SUCCESS: All tickets were booked correctly!")
	} else if stats.BookedTickets < stats.TotalTickets {
		log.Printf("⚠️  WARNING: Only %d out of %d tickets were booked", stats.BookedTickets, stats.TotalTickets)
	}

	if atomic.LoadInt32(&successfulBookings) != stats.BookedTickets {
		log.Printf("❌ ERROR: Booking count mismatch! Counted: %d, Actual: %d", atomic.LoadInt32(&successfulBookings), stats.BookedTickets)
	}
}
