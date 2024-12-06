package handlers

import (
	"auth-service/service"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

func RefreshHandler(authService service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		refreshToken := c.Cookies("refresh_token", "")

		refreshResult, err := authService.RefreshToken(refreshToken)
		if err != nil {
			c.Status(http.StatusForbidden)
			return err
		}

		c.Cookie(&fiber.Cookie{
			Name:  "refresh_token",
			Value: refreshResult.RefreshToken,
		})

		return c.JSON(fiber.Map{
			"access_token": refreshResult.AccessToken,
		})
	}
}
