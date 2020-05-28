package grpc

import (
	grpcDBClient "inventory-go/applications/inventory-db/service/grpc/routes"

	"github.com/caarlos0/env"
	"google.golang.org/grpc"
)

type grpcClientConfig struct {
	ServerAddr string `env:"GRPC_ADDRESS,required"`
}

// Client is a client keeping track of open connexion to RPC server
type Client struct {
	conn              *grpc.ClientConn
	reservationClient *grpcDBClient.ReservationClient
}

// Start opens a new RPC connexion to ServerAddr
func Start() (client *Client, err error) {
	// Retreive Env config
	config := grpcClientConfig{}
	err = env.Parse(&config)
	if err != nil {
		return
	}
	conn, err := grpc.Dial(config.ServerAddr, grpc.WithInsecure())
	if err != nil {
		return
	}

	//create reservation client to db service
	reservationClient := grpcDBClient.NewReservationClient(conn)

	client = &Client{
		conn:              conn,
		reservationClient: &reservationClient,
	}
	return
}

// Close closes the open connexion
func (c *Client) Close() {
	c.conn.Close()
}

// GetReservationClient returns the RPC dedicated Reservation client
func (c *Client) GetReservationClient() grpcDBClient.ReservationClient {
	return *c.reservationClient
}
