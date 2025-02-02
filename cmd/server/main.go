package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	_ "github.com/guemidiborhane/factorydigitale.tech/docs"
	"github.com/guemidiborhane/factorydigitale.tech/internal/config"
	"github.com/guemidiborhane/factorydigitale.tech/internal/http"
	"github.com/guemidiborhane/factorydigitale.tech/internal/logger"
	"github.com/guemidiborhane/factorydigitale.tech/internal/monitor"
	"github.com/guemidiborhane/factorydigitale.tech/internal/router"
	"github.com/guemidiborhane/factorydigitale.tech/internal/storage"
	"github.com/guemidiborhane/factorydigitale.tech/internal/validation"
	websocket "github.com/guemidiborhane/factorydigitale.tech/internal/websocket/server"
	"github.com/guemidiborhane/factorydigitale.tech/pkg"
	"github.com/guemidiborhane/factorydigitale.tech/static"
)

//	@title			FactoryDigitale
//	@version		1.0
//	@description	Interview
//	@contact.name	API
//	@contact.email	guemidiborhane@gmail.com
//	@BasePath		/
func main() {
	app, err := config.NewApp()
	if err != nil {
		log.Fatal(err)
	}

	if config.AppConfig.IsDev() {
		app.Get("/swagger/*", swagger.HandlerDefault)
	}
	logger.Setup()
	Setup(app)
	Start(app)
}

func Setup(a *fiber.App) {
	storage.Setup()
	validation.Setup()

	websocketRedisClient := storage.RedisStorage("websocket").Conn()
	hub := websocket.NewHub(websocketRedisClient)
	router.Setup(a)
	hub.RegisterRoutes(a)
	pkg.Setup(a, hub)

	router.API.Get("/wstest/:channel/:message", hub.TestHandler)

	monitor.Setup(router.API)
	static.RegisterHandler(a)
}

func Start(a *fiber.App) {
	if config.AppConfig.IsDev() {
		storage.WriteRedisEnvVar()
		http.StartServer(a)
	} else {
		http.StartServerWithGracefulShutdown(a)
	}
}
