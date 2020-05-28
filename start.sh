#!/bin/sh

export GO111MODULE=on  # Enable module mode

# Code Generation
protoc ./applications/inventory-db/service/grpc/routes/reservation.proto --go_out=plugins=grpc:.

# Container creation
docker-compose up -V  --build --force-recreate --remove-orphans --always-recreate-deps 
