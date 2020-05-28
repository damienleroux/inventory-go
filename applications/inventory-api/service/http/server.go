package http

import (
	routes "inventory-go/applications/inventory-api/service/http/routes"

	"github.com/caarlos0/env"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type appConfig struct {
	ServerAddress string `env:"HTTP_ADDRESS,required"`
}

// Server is the http server
type Server struct {
	e *echo.Echo
}

// New server
func New(usecase routes.ReservationUseCase) *Server {
	// Create echo server
	server := Server{echo.New()}

	// Log traffic
	server.e.Use(middleware.Logger())

	// Catch and format error
	server.e.Use(ErrorMiddleware)

	// Panic handler
	server.e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize: 1 << 10, // 1 KB
	}))

	// Create http route handler
	CreateReservationHandler := routes.CreateReservation(usecase)

	// Apply handle to route
	server.e.POST("/reservations", CreateReservationHandler)

	return &server
}

// Start server
func (server *Server) Start() (err error) {

	// Retreive Env config
	config := appConfig{}
	err = env.Parse(&config)
	if err != nil {
		return
	}
	return server.e.Start(config.ServerAddress)
}
