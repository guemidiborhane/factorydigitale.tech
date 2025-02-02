package users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/guemidiborhane/factorydigitale.tech/pkg/permissions"
	"github.com/guemidiborhane/factorydigitale.tech/pkg/users/auth"
)

func CheckAuthenticated(c *fiber.Ctx) error {
	if err := auth.VerifyUserSession(c); err != nil {
		return err
	}

	return c.Next()
}

func Can(perms ...string) fiber.Handler {
	middlewares := []fiber.Handler{auth.VerifyUserSession, permissions.CheckPermission(perms...)}

	return func(c *fiber.Ctx) error {
		var err error
		for _, m := range middlewares {
			if err = m(c); err != nil {
				return err
			}
		}

		return err
	}
}
