package handler

import (
	"golang-bootcamp-July/internal/domain"
	"golang-bootcamp-July/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BookingHandler struct {
	bookingService *service.BookingService
}

func NewBookingHandler(bookingService *service.BookingService) *BookingHandler {
	return &BookingHandler{
		bookingService: bookingService,
	}
}

// POST /api/v1/bookings
func (h *BookingHandler) BookTicket(c *gin.Context) {
	var req domain.BookingRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request body",
		})
		return
	}

	result, err := h.bookingService.BookTicketService(c.Request.Context(), req.UserID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == domain.ErrNoTicketsAvailable {
			statusCode = http.StatusConflict
		} else if err == domain.ErrInvalidUserID {
			statusCode = http.StatusBadRequest
		}

		c.JSON(statusCode, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// GET /api/v1/stats
func (h *BookingHandler) GetStats(c *gin.Context) {
	stats, err := h.bookingService.GetBookingStatsService(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
	})
}

// GET /api/v1/bookings/user/:userID
func (h *BookingHandler) GetUserBookigs(c *gin.Context) {
	userID := c.Param("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "User ID is requires",
		})
		return
	}

	bookings, err := h.bookingService.GetUserBookingsService(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    bookings,
	})
}

func (h *BookingHandler) HeathCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    gin.H{"status": "healthy"},
	})
}
