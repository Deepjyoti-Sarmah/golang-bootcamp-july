package main

import (
	"context"
	"golang-bootcamp-July/internal/handler"
	"golang-bootcamp-July/internal/repository/inmemory"
	"golang-bootcamp-July/internal/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// initialize repositories
	bookingRepo := inmemory.NewBookingRepository()
	ticketRepo := inmemory.NewTicketRepository(50_000)

	// initialize service
	bookingService := service.NewBookingService(bookingRepo, ticketRepo)

	// initialize handler
	bookingHandler := handler.NewBookingHandler(bookingService)

	// Setup Gin router
	r := gin.Default()

	// healthCheck
	r.GET("/health", bookingHandler.HeathCheck)

	// API routes
	v1 := r.Group("/api/v1")
	{
		v1.POST("/bookings", bookingHandler.BookTicket)
		v1.GET("/stats", bookingHandler.GetStats)
		v1.GET("/bookings/user/:userID", bookingHandler.GetUserBookigs)
	}

	// Server config
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		log.Println("🚀 Server starting on :8080")
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Gracefull shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Server shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("✅ Server exited")
}
