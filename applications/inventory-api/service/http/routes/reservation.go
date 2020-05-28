package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// OrderStatus is an enum for available reservation status
type OrderStatus string

const (
	//ReservationReserved means the reservation is completed
	ReservationReserved OrderStatus = "RESERVED"
	//ReservationBackOrder means the reservation was completed but refund (NOT USED FOR NOW)
	ReservationBackOrder = "BACKORDER"
	//ReservationPending means the reservation is still processed (NOT USED FOR NOW)
	ReservationPending = "PENDING"
)

// OrderLine is a reservation order line
type OrderLine struct {
	ProductID       string `json:"product"`
	ProductQuantity uint32 `json:"quantity"`
}

// ReservationRequest is the reservation payload request
type ReservationRequest struct {
	Lines []OrderLine `json:"lines"`
}

// ReservationResponse is the reservation payload response
type ReservationResponse struct {
	ID        string      `json:"id"`
	CreatedAt string      `json:"created_at"`
	Lines     []OrderLine `json:"lines"`
	Status    OrderStatus `json:"status"`
}

// ReservationUseCase is the interface that should be implemented to defined was the handler are exepected to do
type ReservationUseCase interface {
	CreateReservation(ReservationRequest) (ReservationResponse, error)
}

// CreateReservation is the route handler for customer submitting new product quantity reservation
func CreateReservation(useCase ReservationUseCase) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		payload := ReservationRequest{}

		if err = c.Bind(&payload); err != nil {
			return
		}

		response, err := useCase.CreateReservation(payload)
		if err != nil {
			return
		}

		return c.JSON(http.StatusCreated, response)
	}

}
