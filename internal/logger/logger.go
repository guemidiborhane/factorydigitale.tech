package logger

import (
	"context"
	"log/slog"
	"os"
	"reflect"
	"time"

	"github.com/guemidiborhane/factorydigitale.tech/internal/config"
	"github.com/gofiber/fiber/v2"
	"github.com/lmittmann/tint"
	slogfiber "github.com/samber/slog-fiber"
)

var (
	output     = os.Stdout
	Logger     *slog.Logger
	Middleware fiber.Handler
)

func Setup() {
	Logger = slog.New(tint.NewHandler(output, &tint.Options{
		AddSource: false,
		Level:     slog.LevelDebug,
	}))

	slog.SetDefault(Logger)
	Logger.With("env", config.AppConfig.Env)
	Middleware = slogfiber.New(Logger)
}

func NewAttribute(key string, value any) slog.Attr {
	var attribute slog.Attr

	switch reflect.TypeOf(value).Kind() {
	case reflect.Int:
		attribute = slog.Int(key, value.(int))
	case reflect.String:
		attribute = slog.String(key, value.(string))
	case reflect.Bool:
		attribute = slog.Bool(key, value.(bool))
	default:
		attribute = slog.Any(key, value)
	}

	if reflect.TypeOf(value) == reflect.TypeOf(time.Now()) {
		attribute = slog.Time(key, value.(time.Time))
	}

	return attribute
}

type Attrs map[string]any

func Log(level slog.Level, msg string, attrs Attrs) {
	var attributes []slog.Attr = make([]slog.Attr, len(attrs))
	for key, value := range attrs {
		attribute := NewAttribute(key, value)

		attributes = append(attributes, attribute)
	}

	attributes = append(attributes, slog.Time("time", time.Now()))

	slog.LogAttrs(
		context.Background(),
		level,
		msg,
		attributes...,
	)
}

func Warn(msg string, attrs Attrs) {
	Log(slog.LevelWarn, msg, attrs)
}

func Error(msg string, attrs Attrs) {
	Log(slog.LevelError, msg, attrs)
}

func Info(msg string, attrs Attrs) {
	Log(slog.LevelInfo, msg, attrs)
}

func Debug(msg string, attrs Attrs) {
	Log(slog.LevelDebug, msg, attrs)
}
