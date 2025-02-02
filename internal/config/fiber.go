package config

import (
	"runtime"

	"github.com/guemidiborhane/factorydigitale.tech/internal/errors"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
)

var FiberConfig = fiber.Config{
	Prefork:           true,
	ReduceMemoryUsage: true,
	Concurrency:       256 * 1024 * 1024,
	StrictRouting:     false,
	CaseSensitive:     false,
	ErrorHandler:      errors.HandleHttpErrors,
	JSONEncoder:       json.Marshal,
	JSONDecoder:       json.Unmarshal,
	// ReadTimeout:   0,
	// WriteTimeout:  0,
	// IdleTimeout:   0,
	// EnableTrustedProxyCheck: false,
	// TrustedProxies:          []string{},
	// EnableIPValidation:      false,
	// EnablePrintRoutes:       true,
}

type c struct {
	Env        string `mapstructure:"APP_ENV"`
	AppKey     string `mapstructure:"APP_KEY"`
	MaxThreads int    `mapstructure:"MAX_THREADS"`
}

var AppConfig = &c{
	Env:        "development",
	AppKey:     "",
	MaxThreads: 2,
}

func NewApp() (*fiber.App, error) {
	if err := EnvFile.LoadConfig(&AppConfig); err != nil {
		return nil, err
	}

	runtime.GOMAXPROCS(AppConfig.MaxThreads)

	return fiber.New(FiberConfig), nil
}

func (c *c) IsDev() bool {
	return c.Env == "development"
}

func (c *c) IsTest() bool {
	return c.Env == "test"
}
