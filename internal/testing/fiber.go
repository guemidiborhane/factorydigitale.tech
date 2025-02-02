package testing

import (
	"github.com/guemidiborhane/factorydigitale.tech/internal/config"
	"github.com/guemidiborhane/factorydigitale.tech/internal/router"
	"github.com/gofiber/fiber/v2"
)

func NewApp() (*fiber.App, error) {
	var a *fiber.App
	a, err := config.NewApp()
	if err != nil {
		return nil, err
	}

	router.Setup(a)

	return a, nil
}
