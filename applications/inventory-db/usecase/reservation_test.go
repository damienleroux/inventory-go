package usecase

import (
	"context"
	grpc "inventory-go/applications/inventory-db/service/grpc/routes"
	"inventory-go/applications/inventory-db/service/postgres"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type FakeReservationClient struct {
	T                        *testing.T
	ReservationID            int
	AvailableProductQuantity map[string]uint32
	CreatedAt                time.Time
}

func (r *FakeReservationClient) Begin() {
}

func (r *FakeReservationClient) Commit() error {
	return nil
}

func (r *FakeReservationClient) Rollback() error {
	return nil
}

func (r *FakeReservationClient) CreateReservation(status string) (result postgres.Reservation, err error) {
	assert.Equal(r.T, status, grpc.ReservationStatus_RESERVED.String(), "When reserving, the reservation status should be ReservationStatus_RESERVED")
	return postgres.Reservation{
		ID:        r.ReservationID,
		Status:    status,
		CreatedAt: r.CreatedAt,
	}, nil
}

func (r *FakeReservationClient) CreateReservationLine(ClientreservationID int, productID string, quantity uint32) (err error) {
	return nil
}

func (r *FakeReservationClient) SelectForUpdateProductAvailibility(productID string) (result postgres.ProductAvailability, err error) {
	return postgres.ProductAvailability{
		Productid: productID,
		Quantity:  r.AvailableProductQuantity[productID],
		UpdatedAt: time.Now(),
	}, nil
}

func (r *FakeReservationClient) UpdateAvailableQuantity(productID string, quantityToReserve uint32) (err error) {
	return nil
}
func TestCreateReservation_getReservableProductQuantity(t *testing.T) {
	assert.Equal(t, getReservableProductQuantity(10, 10), uint32(10), "Take all remaning quantity")
	assert.Equal(t, getReservableProductQuantity(10, 11), uint32(10), "Take more than the remaning quantity")
	assert.Equal(t, getReservableProductQuantity(10, 9), uint32(9), "Take less than the remaning quantity")
}

func TestCreateReservation_ReserveExactAvailableQuantity(t *testing.T) {
	availableProductQuantity := map[string]uint32{
		"PIPR-JACKET-SIZM": 1,
		"PIPR-JACKET-SIZL": 1,
	}

	pgReservationClient := &FakeReservationClient{
		T:                        t,
		ReservationID:            1,
		AvailableProductQuantity: availableProductQuantity,
		CreatedAt:                time.Now(),
	}

	// Create reservation usecase
	reservationUseCase := New(pgReservationClient)

	fakeReq := &grpc.ReserveRequest{
		Lines: []*grpc.ReservationLine{
			{
				Product:  "PIPR-JACKET-SIZM",
				Quantity: 1,
			},
			{
				Product:  "PIPR-JACKET-SIZL",
				Quantity: 1,
			},
		},
	}

	result, err := reservationUseCase.Reserve(context.Background(), fakeReq)

	assert.Equal(t, err, nil, "No error should occurs")
	assert.Equal(t, result.ID, strconv.Itoa(pgReservationClient.ReservationID))
	assert.Equal(t, result.Status, grpc.ReservationStatus_RESERVED)
	assert.Equal(t, result.Lines[0].Product, "PIPR-JACKET-SIZM")
	assert.Equal(t, result.Lines[0].Quantity, uint32(1))
	assert.Equal(t, result.Lines[1].Product, "PIPR-JACKET-SIZL")
	assert.Equal(t, result.Lines[1].Quantity, uint32(1))
	assert.Equal(t, result.CreatedAt, pgReservationClient.CreatedAt.String())
	assert.Nil(t, err)
}

func TestCreateReservation_ReserveLessThanAvailableQuantity(t *testing.T) {
	availableProductQuantity := map[string]uint32{
		"PIPR-JACKET-SIZM": 10,
		"PIPR-JACKET-SIZL": 1,
	}

	pgReservationClient := &FakeReservationClient{
		T:                        t,
		ReservationID:            1,
		AvailableProductQuantity: availableProductQuantity,
		CreatedAt:                time.Now(),
	}

	// Create reservation usecase
	reservationUseCase := New(pgReservationClient)

	fakeReq := &grpc.ReserveRequest{
		Lines: []*grpc.ReservationLine{
			{
				Product:  "PIPR-JACKET-SIZM",
				Quantity: 1,
			},
			{
				Product:  "PIPR-JACKET-SIZL",
				Quantity: 1,
			},
		},
	}

	result, err := reservationUseCase.Reserve(context.Background(), fakeReq)

	assert.Equal(t, err, nil, "No error should occurs")
	assert.Equal(t, result.ID, strconv.Itoa(pgReservationClient.ReservationID))
	assert.Equal(t, result.Status, grpc.ReservationStatus_RESERVED)
	assert.Equal(t, result.Lines[0].Product, "PIPR-JACKET-SIZM")
	assert.Equal(t, result.Lines[0].Quantity, uint32(1))
	assert.Equal(t, result.Lines[1].Product, "PIPR-JACKET-SIZL")
	assert.Equal(t, result.Lines[1].Quantity, uint32(1))
	assert.Equal(t, result.CreatedAt, pgReservationClient.CreatedAt.String())
	assert.Nil(t, err)
}

func TestCreateReservation_ReserveMoreThanAvailableQuantity(t *testing.T) {
	availableProductQuantity := map[string]uint32{
		"PIPR-JACKET-SIZM": 1,
		"PIPR-JACKET-SIZL": 1,
	}

	pgReservationClient := &FakeReservationClient{
		T:                        t,
		ReservationID:            1,
		AvailableProductQuantity: availableProductQuantity,
		CreatedAt:                time.Now(),
	}

	// Create reservation usecase
	reservationUseCase := New(pgReservationClient)

	fakeReq := &grpc.ReserveRequest{
		Lines: []*grpc.ReservationLine{
			{
				Product:  "PIPR-JACKET-SIZM",
				Quantity: 10,
			},
			{
				Product:  "PIPR-JACKET-SIZL",
				Quantity: 1,
			},
		},
	}

	result, err := reservationUseCase.Reserve(context.Background(), fakeReq)

	assert.Equal(t, err, nil, "No error should occurs")
	assert.Equal(t, result.ID, strconv.Itoa(pgReservationClient.ReservationID))
	assert.Equal(t, result.Status, grpc.ReservationStatus_RESERVED)
	assert.Equal(t, result.Lines[0].Product, "PIPR-JACKET-SIZM")
	assert.Equal(t, result.Lines[0].Quantity, uint32(1))
	assert.Equal(t, result.Lines[1].Product, "PIPR-JACKET-SIZL")
	assert.Equal(t, result.Lines[1].Quantity, uint32(1))
	assert.Equal(t, result.CreatedAt, pgReservationClient.CreatedAt.String())
	assert.Nil(t, err)
}
