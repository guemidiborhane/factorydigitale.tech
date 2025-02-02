package users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/guemidiborhane/factorydigitale.tech/pkg/permissions"
	"github.com/guemidiborhane/factorydigitale.tech/pkg/users/auth"
)

func SetupRoutes(r fiber.Router) {
	group := r.Group("/auth")
	group.Post("/users", Can("users:create"), validateRegister, Register)

	authRoutes(group)
	permissionsRoutes(r)
	userRoutes(r)
}

func authRoutes(r fiber.Router) {
	r.Post("/", auth.ValidateLogin, auth.Login)
	r.Delete("/", CheckAuthenticated, auth.Logout)
}

func permissionsRoutes(r fiber.Router) {
	r.Get("/permissions", Can("permissions:index"), permissions.Index)
	r.Put("/permissions/:role", Can("permissions:edit"), permissions.Store)
}

func userRoutes(r fiber.Router) {
	group := r.Group("/user")
	group.Get("/", CheckAuthenticated, Show)
	group.Get("/check", Check)
}
