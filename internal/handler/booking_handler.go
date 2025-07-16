package handler

import (
	"encoding/json"
	"golang-bootcamp-July/internal/domain"
	"golang-bootcamp-July/internal/service"
	"net/http"
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
		return
	}

	result, err := h.bookingService.BookTicketService(r.Context(), req.UserID)
	if err != nil {
		stautsCode := http.StatusInternalServerError
		if err == domain.ErrNoTicketsAvailable {
			stautsCode = http.StatusConflict
		} else if err == domain.ErrInvalidUserID {
			stautsCode = http.StatusBadRequest
		}
		h.sendErrorResponse(w, stautsCode, err.Error())
		return
	}

	h.sendSuccessResponse(w, result)
}

// GET /api/v1/stats
func (h *BookingHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.bookingService.GetBookingStatsService(r.Context())
	if err != nil {
		h.sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
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
		h.sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccessResponse(w, bookings)
}

func (h *BookingHandler) HeathCheck(w http.ResponseWriter, r *http.Request) {
	h.sendSuccessResponse(w, map[string]string{"status": "healthy"})
}

func (h *BookingHandler) sendSuccessResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    data,
	})
}

func (h *BookingHandler) sendErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": false,
		"data":    message,
	})
}
