package handlers

import (
	"auth-service/domain"
	"auth-service/service"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

func SignupHandler(registrationService service.Registration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var registrationUser domain.LoginUser
		if err := c.BodyParser(&registrationUser); err != nil {
			return c.Status(http.StatusForbidden).JSON(err.Error())
		}

		loginResult, err := registrationService.SignUp(registrationUser.Username, registrationUser.Password)
		if err != nil {
			return c.Status(http.StatusForbidden).JSON(err.Error())
		}

		return c.JSON(fiber.Map{
			"access_token":  loginResult.AccessToken,
			"refresh_token": loginResult.RefreshToken,
		})
	}
}
