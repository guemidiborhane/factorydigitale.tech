package static

import (
	"embed"
	"fmt"
	"log"
	"net/http"

	app "github.com/guemidiborhane/factorydigitale.tech/internal/config"
	"github.com/guemidiborhane/factorydigitale.tech/internal/router"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/proxy"
)

//go:embed all:build
var static embed.FS

var FS = filesystem.New(filesystem.Config{
	Root:         http.FS(static),
	Browse:       false,
	PathPrefix:   "/build",
	Index:        "index.html",
	NotFoundFile: "/build/index.html",
})

func RegisterHandler(a *fiber.App) {
	if app.AppConfig.IsDev() {
		log.Println("Running in dev mode")
		setupDevProxy(a)
		return
	}

	a.Use(
		router.CompressMiddleware,
		router.RecoverMiddleware,
		FS,
	)
}

type ViteDevConfig struct {
	Host string `mapstructure:"VITE_HOST"`
	Port int    `mapstructure:"VITE_PORT"`
}

var config = &ViteDevConfig{
	Host: "localhost",
	Port: 5173,
}

func setupDevProxy(a *fiber.App) {
	if err := app.EnvFile.LoadConfig(&config); err != nil {
		log.Fatal(err)
	}

	url := fmt.Sprintf("http://%s:%d", config.Host, config.Port)
	a.Use(
		proxy.Balancer(proxy.Config{
			Servers: []string{url},
			Next: func(c *fiber.Ctx) bool {
				return len(c.Path()) >= 4 && c.Path()[:4] == "/api"
			},
		}),
	)
}
