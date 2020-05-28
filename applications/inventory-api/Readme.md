# inventory-api

## Structure

Two main folders:
- `service`: contains the `HTTP` layer for outside connexions + `gRPC` layer for the connexion with `inventory-db`.
  - `http`: uses [`echo`](https://github.com/labstack/echo) to declare a new server and defines the available routes. The server is able to handle errors gracefully and to log HTTP traffic.
  - `grpc`: uses [`gRPC`](https://grpc.io/docs/guides/) to instanciate the client connexion to the `inventory-db` RPC server. This layer is based on [the gRPC client definition](../inventory-db/service/grpc/routes/reservation.pb.go#L250) declared within `inventory-db` (with `protoc`)
- `usecase`: contains the logic definition to trigger the right RPC calls with the system. This layer requires to instanciate a grpc interface for communication. The struct that defines the interface is provided by protoc generated code within `inventory-db`.

## Test

For now, only one HTTP connexion is tested. That's not an End to End test. In the future, other tests could be added if a whole coverage is expected.

## Caveat

- The API should be based on a declarative file to ease documentation and facilitate maintenance
- The Error code should be declared in a separated file
- The incoming payload are not yet validated, for example using a json schema validator.