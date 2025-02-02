package http

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/guemidiborhane/factorydigitale.tech/internal/config"
	"github.com/guemidiborhane/factorydigitale.tech/internal/storage"
	"github.com/guemidiborhane/factorydigitale.tech/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type serverConfig struct {
	Host string `mapstructure:"HOST"`
	Port uint   `mapstructure:"PORT"`
}

var Config = &serverConfig{
	Host: "0.0.0.0",
	Port: 3000,
}

func StartServerWithGracefulShutdown(a *fiber.App) {
	if err := config.EnvFile.LoadConfig(&Config); err != nil {
		utils.WriteToStderr(err)
	}

	idleConnsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt) // Catch OS signals.
		<-sigint

		if err := a.Shutdown(); err != nil {
			log.Printf("Oops... Server is not shutting down! Reason: %v", err)
		}

		sqlDb, _ := storage.DB.DB()

		defer sqlDb.Close()

		close(idleConnsClosed)
	}()

	url := fmt.Sprintf("%s:%d", Config.Host, Config.Port)

	if err := a.Listen(url); err != nil {
		log.Printf("Oops... Server is not running! Reason: %v", err)
	}

	<-idleConnsClosed
}

func StartServer(a *fiber.App) {
	if err := config.EnvFile.LoadConfig(&Config); err != nil {
		utils.WriteToStderr(err)
	}

	url := fmt.Sprintf("%s:%d", Config.Host, Config.Port)

	if err := a.Listen(url); err != nil {
		log.Printf("Oops... Server is not running! Reason: %v", err)
	}
}
