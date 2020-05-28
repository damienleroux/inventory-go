package postgres

import "time"

//ProductAvailability represents a table product_inventory_availability row
type ProductAvailability struct {
	Productid string    `db:"product_id"`
	Quantity  uint32    `db:"quantity"`
	UpdatedAt time.Time `db:"updated_at"`
}

//Reservation represents a table reservation row
type Reservation struct {
	ID        int       `db:"id"`
	Status    string    `db:"status" pg:"type:reservation_status"`
	CreatedAt time.Time `db:"created_at"`
}
