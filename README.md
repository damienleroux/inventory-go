# Inventory Go

*Disclaimer: This repository contains 2 applications working together to create inventory reservations in a simulated environnement. My intents was to find the right balance between [clean archi principle](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) and [KISS principle](https://en.wikipedia.org/wiki/KISS_principle). How to be relevant in creating Go services without overengineering it?*

Warning: This repository is a test and is not production ready.

## Requirements

Please, make sure your environnement is set with:

- [Docker](https://docs.docker.com/engine/install/)
- [Docker-compose](https://docs.docker.com/compose/install/)
- [Go](https://golang.org/)
- [grpC: protoc & protoc-gen-go](https://grpc.io/docs/quickstart/go/)

## Commands

### Start

Run `sh ./start.sh` to deploy Docker containers

### Test

Run `sh ./start.sh` to run unit test

## Structure

This repository contains 2 applications:
- [`inventory-db`](./applications/inventory-db/Readme.md): a service dedicated to Postgres DB inventory management.
- [`inventory-api`](./applications/inventory-api/Readme.md): a service dedicated to api definition and HTTP connexion.

*Caveat: To keep it simple, business rules, like [quantity computation rules](./applications/inventory-db/usecase/reservation.go#L47), are computed within `inventory-db`. If this application grows and needs complexe & performance consuming business rules, these business rules could be moved to a dedicated service.*

## To go further

Possible future improvements:
- `routes`, through `gRPC` or `HTTP`, are included withing the layer `service`. For the sake of clarity, and to separate `transport methods` from `declarative apis`, this could be split into separated folders.
- Add api routes to list existing registrations or update the available quantity for a dedicated Product. For now [products][./applications/inventory-db/services/postgres/init.sql#L49] are not updatable.
- Solve [`inventory-db` caveats](./applications/inventory-db/Readme.md)
- Solve [`inventory-api` caveats](./applications/inventory-api/Readme.md)
- [docker-compose.yml](./docker-compose.yml) should not instanciate the DB. A dedicated environnement should be available for database interactions, even locally. For now, to ease testing, the DB creation is provided through [docker-compose.yml](./docker-compose.yml).
- We could replace gRPC connexions with message queues like rabbitMQ. This could be useful if asynchronous actions are necessary.
- Add a dictionnary for API error calls
- Build the services and deploy only binaries. "`go run`" should be kept for development purpose only.