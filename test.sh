#!/bin/sh

# Code Generation
protoc ./applications/inventory-db/service/grpc/routes/reservation.proto --go_out=plugins=grpc:.

# run unit tests
go test -v ./...

