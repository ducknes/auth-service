package api

import (
	"auth-service/api/handlers"
	"auth-service/service"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type Server struct {
	server *fiber.App
}

func NewServer() *Server {
	app := fiber.New()
	app.Use(cors.New())

	return &Server{server: app}
}

func (s *Server) Start(port int) {
	if err := s.server.Listen(fmt.Sprintf(":%d", port)); err != nil {
		panic(err)
	}
}

func (s *Server) Setup(authService service.AuthService, registrationService service.Registration) {
	s.server.Post("/sign-up", handlers.SignupHandler(registrationService))
	s.server.Post("/login", handlers.LoginHandler(authService))
	s.server.Post("/logout", handlers.LogoutHandler(authService))
	s.server.Post("/refresh-token", handlers.RefreshHandler(authService))
}
