package reservation

import (
	"errors"
	"net/http"
	"spotsync/internal/domain/reservation/dto"
	"spotsync/internal/httpresponse"
	"strconv"

	"github.com/labstack/echo/v5"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateReservation(c *echo.Context) error {
	userID, ok := c.Get("user_id").(uint)
	if !ok {
		return httpresponse.SendUnauthorized(c, "Unauthorized: user ID not found")
	}

	var req dto.CreateReservationRequest
	if err := c.Bind(&req); err != nil {
		return httpresponse.SendBadRequest(c, "Invalid request payload", err.Error())
	}
	if err := c.Validate(&req); err != nil {
		return httpresponse.SendBadRequest(c, "Validation failed", err.Error())
	}

	response, err := h.service.CreateReservation(userID, req)
	if err != nil {
		if errors.Is(err, ErrZoneFull) {
			return httpresponse.SendConflict(c, "Parking zone is full", err.Error())
		}
		if errors.Is(err, ErrZoneNotFound) {
			return httpresponse.SendNotFound(c, "Parking zone not found")
		}
		return httpresponse.SendInternalServerError(c, "Failed to create reservation")
	}

	return httpresponse.SendSuccess(c, http.StatusCreated, "Reservation confirmed successfully", response)
}

func (h *Handler) GetMyReservations(c *echo.Context) error {
	userID, ok := c.Get("user_id").(uint)
	if !ok {
		return httpresponse.SendUnauthorized(c, "Unauthorized: user ID not found")
	}

	responses, err := h.service.GetMyReservations(userID)
	if err != nil {
		return httpresponse.SendInternalServerError(c, "Failed to retrieve reservations")
	}

	return httpresponse.SendSuccess(c, http.StatusOK, "My reservations retrieved successfully", responses)
}

func (h *Handler) CancelReservation(c *echo.Context) error {
	userID, ok := c.Get("user_id").(uint)
	if !ok {
		return httpresponse.SendUnauthorized(c, "Unauthorized: user ID not found")
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return httpresponse.SendBadRequest(c, "Invalid reservation ID", err.Error())
	}

	if err := h.service.CancelReservation(userID, uint(id)); err != nil {
		if errors.Is(err, ErrReservationNotFound) {
			return httpresponse.SendNotFound(c, "Reservation not found")
		}
		if errors.Is(err, ErrNotOwner) {
			return httpresponse.SendForbidden(c, "You can only cancel your own reservations")
		}
		return httpresponse.SendInternalServerError(c, "Failed to cancel reservation")
	}

	return httpresponse.SendSuccess(c, http.StatusOK, "Reservation cancelled successfully", nil)
}

func (h *Handler) GetAllReservations(c *echo.Context) error {
	responses, err := h.service.GetAllReservations()
	if err != nil {
		return httpresponse.SendInternalServerError(c, "Failed to retrieve reservations")
	}

	return httpresponse.SendSuccess(c, http.StatusOK, "Reservations retrieved successfully", responses)
}
