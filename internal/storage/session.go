package storage

import (
	"time"

	"github.com/guemidiborhane/factorydigitale.tech/internal/utils"
	"github.com/gofiber/fiber/v2/middleware/session"
)

var Session *session.Store

func SetupSession() {
	Session = session.New(session.Config{
		Storage:        RedisStorage("sessions"),
		CookieHTTPOnly: true,
		Expiration:     24 * time.Hour,
		KeyLookup:      "cookie:session_id",
		KeyGenerator:   utils.RandomID,
	})
}
