package usecase

import (
	"context"
	"fmt"
	"strconv"

	grpc "inventory-go/applications/inventory-db/service/grpc/routes"
	"inventory-go/applications/inventory-db/service/postgres"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type pgReservationClient interface {
	Begin()
	Commit() error
	Rollback() error
	CreateReservation(status string) (result postgres.Reservation, err error)
	CreateReservationLine(ClientreservationID int, productID string, quantity uint32) (err error)
	SelectForUpdateProductAvailibility(productID string) (result postgres.ProductAvailability, err error)
	UpdateAvailableQuantity(productID string, quantityToReserve uint32) (err error)
}

// ReservationUseCase is the use case
type ReservationUseCase struct {
	pg pgReservationClient
}

// New creates a new use case
func New(pg pgReservationClient) *ReservationUseCase {
	return &ReservationUseCase{
		pg: pg,
	}
}

// Internal business logic compute available quantity
func getReservableProductQuantity(availableQuantity uint32, requestedQuantity uint32) uint32 {
	if requestedQuantity <= availableQuantity {
		return requestedQuantity
	}
	return availableQuantity
}

// Reserve creates save into DB a reservation
// Reserve uses as input and output grpc models for the sake of simplicity. Custom structs could be used in the future if specific needs are identified
func (r *ReservationUseCase) Reserve(ctx context.Context, req *grpc.ReserveRequest) (result *grpc.ReserveResponse, err error) {
	r.pg.Begin()
	defer func() {
		if err != nil {
			r.pg.Rollback()
			return
		}
		err = r.pg.Commit()
	}()

	reservationStatus := grpc.ReservationStatus_RESERVED
	reservation, err := r.pg.CreateReservation(reservationStatus.String())
	if err != nil {
		return
	}

	var lines []*grpc.ReservationLine
	for _, line := range req.Lines {
		// Lock row to avoid concurrent updates

		productAvailability, lineErr := r.pg.SelectForUpdateProductAvailibility(line.Product)
		if lineErr != nil {
			err = lineErr
			break
		}

		if productAvailability.Productid == "" {
			err = status.Errorf(codes.NotFound, fmt.Sprintf("Product %s not found", line.Product))
			break
		}

		quantityToReserve := getReservableProductQuantity(productAvailability.Quantity, line.Quantity)

		err = r.pg.CreateReservationLine(reservation.ID, line.Product, quantityToReserve)
		if err != nil {
			break
		}

		err = r.pg.UpdateAvailableQuantity(line.Product, quantityToReserve)
		if err != nil {
			break
		}

		lines = append(lines, &grpc.ReservationLine{
			Product:  line.Product,
			Quantity: quantityToReserve,
		})
	}

	if err != nil {
		return
	}

	response := &grpc.ReserveResponse{
		ID:        strconv.Itoa(reservation.ID),
		Status:    reservationStatus,
		Lines:     lines,
		CreatedAt: reservation.CreatedAt.String(),
	}
	return response, nil
}
