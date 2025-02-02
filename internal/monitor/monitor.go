package monitor

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/guemidiborhane/factorydigitale.tech/pkg/permissions"
	"github.com/guemidiborhane/factorydigitale.tech/pkg/users"
)

func Setup(r fiber.Router) {
	permissions.AddPolicy("monitor", "index")
	r.Use("/metrics", users.Can("monitor:index"), monitor.New())
}
