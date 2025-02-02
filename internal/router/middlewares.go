package router

import (
	"time"

	"github.com/guemidiborhane/factorydigitale.tech/internal/config"
	"github.com/guemidiborhane/factorydigitale.tech/internal/errors"
	"github.com/guemidiborhane/factorydigitale.tech/internal/storage"
	"github.com/guemidiborhane/factorydigitale.tech/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

var (
	HelmetMiddleware   = helmet.New()
	CorsMiddleware     = cors.New(cors.Config{})
	RecoverMiddleware  = recover.New()
	CompressMiddleware = compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	})
	EtagMiddleware        = etag.New()
	IdempotencyMiddleware = idempotency.New()
	RequestIDMiddleware   = requestid.New()
)

func CsrfMiddleware() fiber.Handler {
	csrfHeaderName := "X-CSRF-Token"

	return csrf.New(csrf.Config{
		KeyLookup:         "header:" + csrfHeaderName,
		CookieName:        "csrf_",
		CookieSameSite:    "Strict",
		CookieSessionOnly: true,
		CookieHTTPOnly:    false,
		Expiration:        1 * time.Hour,
		Extractor:         csrf.CsrfFromHeader(csrfHeaderName),
		Session:           storage.Session,
		SessionKey:        "fiber.csrf.token",
		HandlerContextKey: "fiber.csrf.handler",
		KeyGenerator:      utils.RandomID,
		ErrorHandler:      errors.HandleHttpErrors,
	})
}

func EncryptCookies() fiber.Handler {
	key := config.AppConfig.AppKey

	if len(key) == 0 {
		panic("APP_KEY not set")
	}

	return encryptcookie.New(encryptcookie.Config{
		Key: key,
	})
}

func Sleep(timeout int) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if !config.AppConfig.IsDev() {
			return c.Next()
		}

		time.Sleep(time.Duration(timeout) * time.Second)
		return c.Next()
	}
}
