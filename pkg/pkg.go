package pkg

import (
	"github.com/gofiber/fiber/v2"
	"github.com/guemidiborhane/factorydigitale.tech/internal/router"
	"github.com/guemidiborhane/factorydigitale.tech/internal/setup"
	"github.com/guemidiborhane/factorydigitale.tech/internal/storage"
	websocket "github.com/guemidiborhane/factorydigitale.tech/internal/websocket/server"
	"github.com/guemidiborhane/factorydigitale.tech/pkg/movies"
	"github.com/guemidiborhane/factorydigitale.tech/pkg/permissions"
	"github.com/guemidiborhane/factorydigitale.tech/pkg/users"
)

func register(c *setup.Config) {
	permissions.Setup(c)
	users.Setup(c)
	movies.Setup(c)
}

func Setup(a *fiber.App, hub *websocket.Hub) {
	register(&setup.Config{
		App:      &a,
		Router:   &router.API,
		Database: &storage.DB,
		Session:  &storage.Session,
	})
}
