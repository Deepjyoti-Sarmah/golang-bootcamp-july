package main

import (
	"context"
	"fmt"
	"golang-bootcamp-July/booking"
	"golang-bootcamp-July/worker"
	"log"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	const (
		totalTickets  = 50000
		totalRequests = 75000
		numWorkers    = 1000
	)

	log.Println("🎫 Ticket Booking System Starting...")
	log.Printf("📊 Total Tickets Available: %d", totalTickets)
	log.Printf("👥 Total Booking Requests: %d", totalRequests)
	log.Printf("⚙️  Worker Pool Size: %d", numWorkers)
	log.Println(strings.Repeat("-", 50))

	startTime := time.Now()

	ticketSystem := booking.NewTicketSystem(totalTickets)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	workerPool := worker.NewWorkerPool(numWorkers, ticketSystem)
	workerPool.StartWorkerPool(ctx)

	var (
		successfulBookings int32
		failedBookings     int32
		processedRequests  int32
	)

	resultsDone := make(chan bool)
	go func() {
		for result := range workerPool.Results() {
			atomic.AddInt32(&processedRequests, 1)

			if result.Success {
				atomic.AddInt32(&successfulBookings, 1)
				if processedRequests%5000 == 0 {
					log.Printf("✅ Progress: %d/%d requests processed, %d tickets booked",
						atomic.LoadInt32(&processedRequests),
						totalRequests,
						atomic.LoadInt32(&successfulBookings))
				}
			} else {
				atomic.AddInt32(&failedBookings, 1)
			}
		}
		resultsDone <- true
	}()

	log.Println("🚀 Starting concurrent booking requests...")

	// client simulation (it has nothing to do with server)
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

	<-resultsDone

	duration := time.Since(startTime)
	bookedTickets, remainingTickets := ticketSystem.GetStats()

	log.Println(strings.Repeat("=", 50))
	log.Println("📈 BOOKING SYSTEM FINAL REPORT")
	log.Println(strings.Repeat("=", 50))
	log.Printf("⏱️  Total Time Taken: %v", duration)
	log.Printf("🎫 Total Tickets: %d", totalTickets)
	log.Printf("✅ Total Tickets Booked: %d", bookedTickets)
	log.Printf("❌ Total Tickets NOT Booked: %d", atomic.LoadInt32(&failedBookings))
	log.Printf("📊 Total Requests Processed: %d", atomic.LoadInt32(&processedRequests))
	log.Printf("🎯 Remaining Tickets: %d", remainingTickets)
	log.Printf("⚡ Requests per Second: %.2f", float64(totalRequests)/duration.Seconds())
	log.Printf("🔄 Average Time per Request: %v", duration/time.Duration(totalRequests))
	log.Println(strings.Repeat("=", 50))

	if bookedTickets == totalTickets && remainingTickets == 0 {
		log.Println("✅ SUCCESS: All tickets were booked correctly!")
	} else if bookedTickets < totalTickets {
		log.Printf("⚠️  WARNING: Only %d out of %d tickets were booked", bookedTickets, totalTickets)
	}

	if atomic.LoadInt32(&successfulBookings) != bookedTickets {
		log.Printf("❌ ERROR: Booking count mismatch! Counted: %d, Actual: %d",
			atomic.LoadInt32(&successfulBookings), bookedTickets)
	}
}
