package users

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/guemidiborhane/factorydigitale.tech/internal/errors"
	"github.com/guemidiborhane/factorydigitale.tech/internal/storage"
	"github.com/guemidiborhane/factorydigitale.tech/pkg/permissions"
	"github.com/guemidiborhane/factorydigitale.tech/pkg/users/auth"
)

// @Summary	Register User
// @Tags		Users
// @Accept		json
// @Produce	json
// @Param		body	body		UserRegisterParams	true	"Register request"
// @Success	200		{object}	auth.UserResponse
// @Failure	403		{object}	errors.HttpError
// @Failure	500		{object}	errors.HttpError
// @Router		/api/users [post]
func Register(c *fiber.Ctx) error {
	var user auth.User

	if err := c.BodyParser(&user); err != nil {
		return errors.Unexpected(err.Error())
	}

	if err := user.Create(); err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(user.AsJSON())
}

type UserJSONResponse struct {
	User        *auth.UserResponse `json:"user" binding:"required"`
	Permissions []string           `json:"permissions" binding:"required"`
	Tracks      auth.UserTracks    `json:"tracks" binding:"required"`
}

// @Summary	Show User
// @Tags		Users
// @Accept		json
// @Produce	json
// @Success	200	{object}	UserJSONResponse
// @Failure	401	{object}	errors.HttpError
// @Failure	403	{object}	errors.HttpError
// @Failure	500	{object}	errors.HttpError
// @Router		/api/user [get]
func Show(c *fiber.Ctx) error {
	user, err := auth.GetCurrentUser(c)
	if err != nil {
		return err
	}
	tracks, err := user.GetTracks()
	if err != nil {
		log.Println(err)
	}

	return c.Status(fiber.StatusOK).JSON(UserJSONResponse{
		User:        user.AsJSON(),
		Permissions: permissions.GetRolePermissions(user.Role),
		Tracks:      tracks,
	})
}

// @Summary	Check
// @Tags		Users
// @Accept		json
// @Produce json
// @Success	200	{boolean} boolean true
// @Failure	401	{boolean} boolean false
// @Router		/api/user/check [get]
func Check(c *fiber.Ctx) error {
	session, err := storage.Session.Get(c)

	if err != nil || session.Get(auth.AUTH_KEY) == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(false)
	}

	return c.Status(fiber.StatusOK).JSON(true)
}
