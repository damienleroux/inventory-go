package main

import (
	"log"

	"inventory-go/applications/inventory-db/service/grpc"
	"inventory-go/applications/inventory-db/service/postgres"
	"inventory-go/applications/inventory-db/usecase"
)

func main() {
	// Opens connexion to Postgres Client
	pgClient, err := postgres.Start()
	if err != nil {
		log.Fatal(err)
	}
	defer pgClient.Close()

	// Creates reservation Postgres Client
	pgReservationClient := postgres.NewReservationClient(pgClient)
	if err != nil {
		log.Fatal(err)
	}

	// Creates reservation usecase
	reservationUseCase := usecase.New(pgReservationClient)

	// Creates app, based on reservation usecase
	app := grpc.New(reservationUseCase)

	// Starts app
	err = app.Start()
	if err != nil {
		log.Fatal(err)
	}
}
