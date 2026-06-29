package parkingzone

import (
	"errors"
	"net/http"
	"spotsync/internal/domain/parkingzone/dto"
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

func (h *Handler) CreateZone(c *echo.Context) error {
	var req dto.CreateZoneRequest
	if err := c.Bind(&req); err != nil {
		return httpresponse.SendBadRequest(c, "Invalid request payload", err.Error())
	}
	if err := c.Validate(&req); err != nil {
		return httpresponse.SendBadRequest(c, "Validation failed", err.Error())
	}

	response, err := h.service.Create(req)
	if err != nil {
		return httpresponse.SendInternalServerError(c, "Failed to create parking zone")
	}

	return httpresponse.SendSuccess(c, http.StatusCreated, "Parking zone created successfully", response)
}

func (h *Handler) GetAllZones(c *echo.Context) error {
	responses, err := h.service.GetAll()
	if err != nil {
		return httpresponse.SendInternalServerError(c, "Failed to retrieve parking zones")
	}
	return httpresponse.SendSuccess(c, http.StatusOK, "Parking zones retrieved successfully", responses)
}

func (h *Handler) GetZoneByID(c *echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return httpresponse.SendBadRequest(c, "Invalid zone ID", err.Error())
	}

	response, err := h.service.GetByID(uint(id))
	if err != nil {
		if errors.Is(err, ErrZoneNotFound) {
			return httpresponse.SendNotFound(c, "Parking zone not found")
		}
		return httpresponse.SendInternalServerError(c, "Failed to retrieve parking zone")
	}

	return httpresponse.SendSuccess(c, http.StatusOK, "Parking zone retrieved successfully", response)
}

func (h *Handler) UpdateZone(c *echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return httpresponse.SendBadRequest(c, "Invalid zone ID", err.Error())
	}

	var req dto.UpdateZoneRequest
	if err := c.Bind(&req); err != nil {
		return httpresponse.SendBadRequest(c, "Invalid request payload", err.Error())
	}
	if err := c.Validate(&req); err != nil {
		return httpresponse.SendBadRequest(c, "Validation failed", err.Error())
	}

	response, err := h.service.Update(uint(id), req)
	if err != nil {
		if errors.Is(err, ErrZoneNotFound) {
			return httpresponse.SendNotFound(c, "Parking zone not found")
		}
		return httpresponse.SendInternalServerError(c, "Failed to update parking zone")
	}

	return httpresponse.SendSuccess(c, http.StatusOK, "Parking zone updated successfully", response)
}

func (h *Handler) DeleteZone(c *echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return httpresponse.SendBadRequest(c, "Invalid zone ID", err.Error())
	}

	if err := h.service.Delete(uint(id)); err != nil {
		if errors.Is(err, ErrZoneNotFound) {
			return httpresponse.SendNotFound(c, "Parking zone not found")
		}
		return httpresponse.SendInternalServerError(c, "Failed to delete parking zone")
	}

	return httpresponse.SendSuccess(c, http.StatusOK, "Parking zone deleted successfully", nil)
}
