package users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/guemidiborhane/factorydigitale.tech/internal/errors"
	"github.com/guemidiborhane/factorydigitale.tech/internal/validation"
)

type UserRegisterParams struct {
	Username          string `json:"username" validate:"required,min=3,max=20" gorm:"uniqueIndex"`
	Password          string `json:"password" validate:"required,min=6"`
	Role              string `json:"role" validate:"required"`
	AuthenticableType string `json:"authenticable_type" validate:"required"`
	AuthenticableID   int    `json:"authenticable_id" validate:"required"`
}

func validateRegister(c *fiber.Ctx) error {
	var body UserRegisterParams

	if err := c.BodyParser(&body); err != nil {
		return errors.Unexpected(err.Error())
	}

	if err := validation.Validate(body); err != nil {
		return err
	}

	return c.Next()
}
