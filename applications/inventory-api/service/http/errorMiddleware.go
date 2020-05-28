package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ErrorResponse is the standard error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// ErrorMiddleware is an echo middleware to handle error and return standard HTTP Error code & payoad
func ErrorMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := next(c); err != nil {
			errStatus, _ := status.FromError(err)
			if codes.NotFound == errStatus.Code() {
				return c.JSON(http.StatusNotFound, ErrorResponse{
					Error: errStatus.Message(),
				})
			}
			return c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: errStatus.Message(),
			})
		}
		return nil
	}
}
