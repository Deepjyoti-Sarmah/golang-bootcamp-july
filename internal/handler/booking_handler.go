package handler

import (
	"encoding/json"
	"golang-bootcamp-July/internal/domain"
	"golang-bootcamp-July/internal/service"
	"net/http"
	"strings"
)

type BookingHandler struct {
	bookingService *service.BookingService
	// logger         *logger.Logger
}

func NewBookingHandler(bookingService *service.BookingService) *BookingHandler {
	return &BookingHandler{
		bookingService: bookingService,
	}
}

// POST /api/v1/bookings
func (h *BookingHandler) BookingTicket(w http.ResponseWriter, r *http.Request) {
	var req domain.BookingRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, "Invalid request body")
	}

	result, err := h.bookingService.BookTicketService(r.Context(), req.UserID)
	if err != nil {
		h.sendErrorResponse(w, http.StatusInternalServerError, err.Error())
	}

	// check for the result and commit
	if !result.Success {
		statusCode := http.StatusInternalServerError
		if strings.Contains(result.Error, "no tickets available") {
			statusCode = http.StatusConflict
		} else if strings.Contains(result.Error, "invalid user ID") {
			statusCode = http.StatusBadRequest
		}
		h.sendErrorResponse(w, statusCode, result.Error)
	} else {
		h.sendSuccessResponse(w, result)
	}
}

// GET /api/v1/stats
func (h *BookingHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.bookingService.GetBookingStatsService(r.Context())
	if err != nil {
		h.sendErrorResponse(w, http.StatusInternalServerError, err.Error())
	}

	h.sendSuccessResponse(w, stats)
}

// GET /api/v1/bookings/user/{userID}
func (h *BookingHandler) GetUserBookings(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userID")
	if userID == "" {
		h.sendErrorResponse(w, http.StatusBadRequest, "User ID is required")
		return
	}

	bookings, err := h.bookingService.BookTicketService(r.Context(), userID)
	if err != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	h.sendSuccessResponse(w, bookings)
}

func (h *BookingHandler) HeathCheck(w http.ResponseWriter, r *http.Request) {}

func (h *BookingHandler) sendSuccessResponse(w http.ResponseWriter, data interface{}) {}

func (h *BookingHandler) sendErrorResponse(w http.ResponseWriter, statusCode int, message string) {}
