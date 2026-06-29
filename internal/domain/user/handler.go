package user

import (
	"errors"
	"net/http"
	"spotsync/internal/domain/user/dto"
	"spotsync/internal/httpresponse"

	"github.com/labstack/echo/v5"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(c *echo.Context) error {
	var req dto.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return httpresponse.SendBadRequest(c, "Invalid request payload", err.Error())
	}
	if err := c.Validate(&req); err != nil {
		return httpresponse.SendBadRequest(c, "Validation failed", err.Error())
	}

	response, err := h.service.Register(req)
	if err != nil {
		if errors.Is(err, ErrorAlreadyExist) {
			return httpresponse.SendConflict(c, "Failed to register user", err.Error())
		}
		return httpresponse.SendInternalServerError(c, "Failed to register user")
	}

	return httpresponse.SendSuccess(c, http.StatusCreated, "User registered successfully", response)
}

func (h *Handler) Login(c *echo.Context) error {
	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return httpresponse.SendBadRequest(c, "Invalid request payload", err.Error())
	}
	if err := c.Validate(&req); err != nil {
		return httpresponse.SendBadRequest(c, "Validation failed", err.Error())
	}

	response, err := h.service.Login(req)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			return httpresponse.SendUnauthorized(c, "Invalid email or password")
		}
		return httpresponse.SendInternalServerError(c, "Failed to login user")
	}

	return httpresponse.SendSuccess(c, http.StatusOK, "Login successful", response)
}
