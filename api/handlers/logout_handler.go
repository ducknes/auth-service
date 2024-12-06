package handlers

import (
	"auth-service/service"
	"github.com/gofiber/fiber/v2"
)

func LogoutHandler(authService service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authService.Logout(c.Cookies("refresh_token", ""))

		c.Cookie(&fiber.Cookie{
			Name:  "refresh_token",
			Value: "dead",
		})
		
		return nil
	}
}
