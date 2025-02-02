package auth

import (
	"github.com/guemidiborhane/factorydigitale.tech/internal/errors"
	"github.com/guemidiborhane/factorydigitale.tech/internal/storage"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

var (
	AUTH_KEY           string = "authenticated"
	USER_ID            string = "user_id"
	AccountLockedError        = &errors.HttpError{Status: 420, Code: "unauthorized", Message: "You're not authorized"}
)

func attemptToLogin(c *fiber.Ctx) (*User, error) {
	var body LoginParams

	if err := c.BodyParser(&body); err != nil {
		return &User{}, errors.BadRequest(err)
	}

	user := &User{
		Username: body.Username,
	}

	if err := user.Get(); err != nil {
		return &User{}, err
	}

	if !user.LockedAt.IsZero() {
		return &User{}, AccountLockedError
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		return &User{}, errors.Unauthorized
	}

	tracks, err := user.GetTracks()
	if err != nil {
		return &User{}, err
	}
	tracks.LoginCount = tracks.LoginCount + 1
	user.WriteTracks(&tracks)

	return user, nil
}

func (user *User) SetSession(c *fiber.Ctx) error {
	session, err := storage.Session.Get(c)
	if err != nil {
		return errors.Unexpected(err.Error())
	}

	session.Set(AUTH_KEY, true)
	session.Set(USER_ID, user.ID)

	if err := session.Save(); err != nil {
		return errors.Unexpected(err.Error())
	}

	return nil
}

// @Summary	Login user
// @Tags		Users
// @Accept		json
// @Produce	json
// @Param		body	body		LoginParams	true	"LoginRequest"
// @Success	200		{object}	UserResponse
// @Failure	401		{object}	errors.HttpError
// @Failure	404		{object}	errors.HttpError
// @Failure	500		{object}	errors.HttpError
// @Router		/api/auth [post]
func Login(c *fiber.Ctx) error {
	user, err := attemptToLogin(c)
	if err != nil {
		return err
	}

	if err := user.SetSession(c); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(user.AsJSON())
}

// @Summary	Logout user
// @Tags		Users
// @Accept		json
// @Produce	json
// @Success	200	{object}	nil
// @Failure	500	{object}	errors.HttpError
// @Router		/api/auth [delete]
func Logout(c *fiber.Ctx) error {
	session, err := storage.Session.Get(c)
	if err != nil {
		return errors.Unexpected(err.Error())
	}

	if err := session.Destroy(); err != nil {
		return errors.Unexpected(err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(nil)
}
