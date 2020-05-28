package grpc

import (
	"net"

	routes "inventory-go/applications/inventory-db/service/grpc/routes"

	"github.com/caarlos0/env"
	"google.golang.org/grpc"
)

type appConfig struct {
	ServerAddr string `env:"GRPC_ADDRESS,required"`
}

// Server discribes grpc server
type Server struct {
	s *grpc.Server
}

// New creates new grpc server
func New(usecase routes.ReservationServer) *Server {
	server := Server{grpc.NewServer()}

	routes.RegisterReservationServer(server.s, usecase)
	return &server
}

// Start starts grpc server
func (server *Server) Start() (err error) {
	// Retreive Env config
	config := appConfig{}
	err = env.Parse(&config)
	if err != nil {
		return
	}

	netlistener, err := net.Listen("tcp", config.ServerAddr)
	if err != nil {
		return
	}

	return server.s.Serve(netlistener)
}
