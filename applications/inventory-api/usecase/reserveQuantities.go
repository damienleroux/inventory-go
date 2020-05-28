package usecase

import (
	"context"
	"errors"
	http "inventory-go/applications/inventory-api/service/http/routes"
	grpc "inventory-go/applications/inventory-db/service/grpc/routes"
)

// ReserveQuantities is a use case implementing routes.ReservationUseCase
// a connexion to gRPC server is request to additional call to inventory-db
type ReserveQuantities struct {
	dbClient grpc.ReservationClient
}

// New creates a new use case
func New(dbClient grpc.ReservationClient) *ReserveQuantities {
	return &ReserveQuantities{dbClient}
}

func marshallDBRequest(orderLines []http.OrderLine) grpc.ReserveRequest {
	request := grpc.ReserveRequest{}

	for _, orderLine := range orderLines {
		requestLine := &grpc.ReservationLine{
			Product:  orderLine.ProductID,
			Quantity: orderLine.ProductQuantity,
		}
		request.Lines = append(request.Lines, requestLine)
	}
	return request
}

func unmarshallDBResponseStatus(responseStatus grpc.ReservationStatus) (status http.OrderStatus, err error) {
	switch responseStatus {
	case grpc.ReservationStatus_BACKORDER:
		status = http.ReservationBackOrder
	case grpc.ReservationStatus_PENDING:
		status = http.ReservationPending
	case grpc.ReservationStatus_RESERVED:
		status = http.ReservationReserved
	default:
		// @todo: create common error handler for the whole app
		err = errors.New("[STATUS ERROR] unknown reservation status")
	}
	return
}

func unmarshallDBResponse(response grpc.ReserveResponse) (reservation http.ReservationResponse, err error) {
	status, err := unmarshallDBResponseStatus(response.Status)
	if err != nil {
		return
	}

	lines := []http.OrderLine{}
	for _, responseLine := range response.Lines {
		orderLine := http.OrderLine{
			ProductID:       responseLine.Product,
			ProductQuantity: responseLine.Quantity,
		}
		lines = append(lines, orderLine)
	}

	reservation = http.ReservationResponse{
		ID:        response.ID,
		Status:    status,
		Lines:     lines,
		CreatedAt: response.CreatedAt,
	}
	return
}

// CreateReservation creates a reservation by calling an external grpc service to ask for db save
// CreateReservation uses as input and output http models for the sake of simplicity.
// Custom structs could be used in the future if specific needs are identified
func (usecase *ReserveQuantities) CreateReservation(reservation http.ReservationRequest) (result http.ReservationResponse, err error) {
	// Prepare request fo db service
	dbRequest := marshallDBRequest(reservation.Lines)

	// Trigger request for DB service and catch response & err
	dbResponse, err := usecase.dbClient.Reserve(context.Background(), &dbRequest)
	if err != nil {
		return
	}
	// Convert DB service response into a standard http response
	result, err = unmarshallDBResponse(*dbResponse)

	return
}
