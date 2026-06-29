package httpresponse

import "net/http"

type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors,omitempty"`
}

func NewSuccess(message string, data interface{}) *SuccessResponse {
	return &SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
}

func NewError(message string, errors interface{}) *ErrorResponse {
	return &ErrorResponse{
		Success: false,
		Message: message,
		Errors:  errors,
	}
}

func SendSuccess(c interface{ JSON(int, interface{}) error }, code int, message string, data interface{}) error {
	return c.JSON(code, NewSuccess(message, data))
}

func SendError(c interface{ JSON(int, interface{}) error }, code int, message string, errors interface{}) error {
	return c.JSON(code, NewError(message, errors))
}

func SendInternalServerError(c interface{ JSON(int, interface{}) error }, message string) error {
	return c.JSON(http.StatusInternalServerError, NewError(message, nil))
}

func SendBadRequest(c interface{ JSON(int, interface{}) error }, message string, errors interface{}) error {
	return c.JSON(http.StatusBadRequest, NewError(message, errors))
}

func SendUnauthorized(c interface{ JSON(int, interface{}) error }, message string) error {
	return c.JSON(http.StatusUnauthorized, NewError(message, nil))
}

func SendForbidden(c interface{ JSON(int, interface{}) error }, message string) error {
	return c.JSON(http.StatusForbidden, NewError(message, nil))
}

func SendNotFound(c interface{ JSON(int, interface{}) error }, message string) error {
	return c.JSON(http.StatusNotFound, NewError(message, nil))
}

func SendConflict(c interface{ JSON(int, interface{}) error }, message string, errors interface{}) error {
	return c.JSON(http.StatusConflict, NewError(message, errors))
}
