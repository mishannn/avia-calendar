package rest

import (
	"github.com/go-playground/validator"
	"github.com/mishannn/avia-calendar/frontend"
	"github.com/mishannn/avia-calendar/internal/services/tickets"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	router         *echo.Echo
	ticketsService *tickets.Service
}

func NewServer(ticketsService *tickets.Service) (*Server, error) {
	server := &Server{
		router:         echo.New(),
		ticketsService: ticketsService,
	}

	server.router.Validator = &RequestValidator{validator: validator.New()}

	server.registerSPA()
	server.registerRoutes()

	return server, nil
}

func (s *Server) registerSPA() {
	s.router.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Filesystem: frontend.BuildHTTPFS(),
		HTML5:      true,
	}))
}

func (s *Server) registerRoutes() {
	s.router.GET("/api/v1/find-cheapest-tickets", s.findCheapestTickets)
	s.router.GET("/api/v1/find-iata", s.findIATA)
}

func (s *Server) Run(addr string) error {
	return s.router.Start(addr)
}
