package auth

import (
	"github.com/guemidiborhane/factorydigitale.tech/internal/errors"
	"github.com/guemidiborhane/factorydigitale.tech/internal/validation"
	"github.com/gofiber/fiber/v2"
)

type LoginParams struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func ValidateLogin(c *fiber.Ctx) error {
	var body LoginParams

	if err := c.BodyParser(&body); err != nil {
		return errors.BadRequest(err.Error())
	}

	if err := validation.Validate(body); err != nil {
		return err
	}

	return c.Next()
}
