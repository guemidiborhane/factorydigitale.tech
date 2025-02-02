package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/guemidiborhane/factorydigitale.tech/internal/logger"
)

var API fiber.Router

func setupAPIRoutes(a *fiber.App) {
	API = a.Group("/api",
		EncryptCookies(),
		HelmetMiddleware,
		CorsMiddleware,
		logger.Middleware,
		CompressMiddleware,
		EtagMiddleware,
		CsrfMiddleware(),
		RequestIDMiddleware,
	)
}
