package movies

import (
	"github.com/guemidiborhane/factorydigitale.tech/internal/setup"
	"github.com/guemidiborhane/factorydigitale.tech/pkg/movies/models"
	"github.com/guemidiborhane/factorydigitale.tech/pkg/permissions"
)

func Setup(c *setup.Config) {
	permissions.RegisterPermissions("movies", permissions.DefaultActions)
	models.Setup(*c.Database)
	setupRoutes((*c.Router))
}
