package setup

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"gorm.io/gorm"
)

type Config struct {
	App      **fiber.App
	Router   *fiber.Router
	Database **gorm.DB
	Session  **session.Store
}
