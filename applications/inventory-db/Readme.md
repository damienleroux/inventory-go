# inventory-db

## Structure

Two main folders:
- `service`: contains the gRPC layer for the connexion from `inventory-api` and the `Postgres` layer to communicate with the `Postgres` DB.
  - `grpc`: defines the RPC services for reservation requests. With `protoc`, a [`client` and a `server`](./service/grpc/routes/reservation.pb.go) are generated. This layer instanciates the RPC `server`, called by `inventory-api` to communicate.
  - `postgres`: offers the DB `client`, defines routes to execute requests to the database and provides the initial database definition.
- `usecase`: contains the business logic for new reservation. This layer requires to instanciate the postgres client interface.

## Test

For now, only the reservation use case is tested to check that the business logic answers expectations.

## Caveat

- The initial database definition is not production ready: it populates testing data for testing purpose. Next step will be:
  - to declare a dedicated file for it
  - to use [DB migrate](https://db-migrate.readthedocs.io/en/latest/) for DB schema management