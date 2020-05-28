package main

import (
	"inventory-go/applications/inventory-api/service/grpc"
	"inventory-go/applications/inventory-api/service/http"
	"inventory-go/applications/inventory-api/usecase"
	"log"
)

func main() {
	// Opens connexion to db service
	grpcConn, err := grpc.Start()
	if err != nil {
		log.Fatal(err)
	}
	defer grpcConn.Close()

	// Creates reservation usecase
	reservationUseCase := usecase.New(grpcConn.GetReservationClient())

	// Creates app, based on reservation usecase
	app := http.New(reservationUseCase)

	// Starts app
	err = app.Start()
	if err != nil {
		log.Fatal(err)
	}
}
