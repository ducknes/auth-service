package handlers

import (
	"auth-service/domain"
	"auth-service/service"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

func LoginHandler(authService service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var loginUser domain.LoginUser
		if err := c.BodyParser(&loginUser); err != nil {
			c.Status(http.StatusForbidden)
			return err
		}

		loginResult, err := authService.Login(loginUser.Username, loginUser.Password)
		if err != nil {
			c.Status(http.StatusForbidden)
			return err
		}

		c.Cookie(&fiber.Cookie{
			Name:  "refresh_token",
			Value: loginResult.RefreshToken,
		})

		return c.JSON(fiber.Map{
			"access_token": loginResult.AccessToken,
		})
	}
}
