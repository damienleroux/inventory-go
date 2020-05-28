package postgres

import (
	"github.com/jmoiron/sqlx"
)

// ReservationClient is a dedicated postgres client for reservation
// ReservationClient implements the use case interface "pgReservationClient"
// The common Postgres client is required to work
type ReservationClient struct {
	pg *Client
	tx *sqlx.Tx
}

// NewReservationClient creates new client
func NewReservationClient(pg *Client) *ReservationClient {
	return &ReservationClient{
		pg: pg,
	}
}

// Begin begins transaction
func (r *ReservationClient) Begin() {
	r.tx = r.pg.Begin()
}

// Commit commits current transaction
func (r *ReservationClient) Commit() error {
	return r.tx.Commit()
}

// Rollback rollback current transaction
func (r *ReservationClient) Rollback() error {
	return r.tx.Rollback()
}

// CreateReservation register a new reservation in data base
func (r *ReservationClient) CreateReservation(status string) (result Reservation, err error) {
	rows, err := r.tx.Queryx(
		`
			INSERT INTO reservation (status)
			VALUES ($1)
			RETURNING id, status, created_at
		`,
		status,
	)
	defer rows.Close()

	result = Reservation{}
	for rows.Next() {
		err = rows.StructScan(&result)
		if err != nil {
			return
		}
	}

	return
}

// CreateReservationLine register a new reservation line in data base
func (r *ReservationClient) CreateReservationLine(ClientreservationID int, productID string, quantity uint32) (err error) {
	rows, err := r.tx.Query(
		`
			INSERT INTO reservation_line (reservation_id, product_id, quantity) 
			VALUES ($1, $2, $3)
		`,
		ClientreservationID,
		productID,
		quantity,
	)
	rows.Close()

	return
}

// SelectForUpdateProductAvailibility locks product quantity to avoid race condition
// This ensures the right quantity management
// Warning: this doesn't prevent DB deadlock. Postgres will canceled the deadlock as soon as detected
// In the future, using either a mutex or using a unique rabbitmq message consumer to update availability will deprecate this row lock
func (r *ReservationClient) SelectForUpdateProductAvailibility(productID string) (result ProductAvailability, err error) {
	rows, err := r.tx.Queryx(
		`
			SELECT product_id, quantity from product_inventory_availability 
			WHERE product_id=$1
			FOR UPDATE;
		`,
		productID,
	)
	defer rows.Close()

	result = ProductAvailability{}

	for rows.Next() {
		err = rows.StructScan(&result)
		if err != nil {
			return
		}
	}

	return
}

// UpdateAvailableQuantity updates the available product quantity
func (r *ReservationClient) UpdateAvailableQuantity(productID string, quantityToReserve uint32) (err error) {
	rows, err := r.tx.Queryx(
		`
			UPDATE product_inventory_availability
			SET quantity = quantity - $2
			WHERE product_id = $1;
		`,
		productID,
		quantityToReserve,
	)
	rows.Close()

	return
}
