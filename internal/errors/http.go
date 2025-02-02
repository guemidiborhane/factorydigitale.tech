package errors

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

func HandleHttpErrors(ctx *fiber.Ctx, err error) error {
	if e, ok := err.(*HttpError); ok {
		return ctx.Status(e.Status).JSON(e)
	} else if e, ok := err.(*fiber.Error); ok {
		return ctx.Status(e.Code).JSON(HttpError{Status: e.Code, Code: "internal-server", Message: e.Message})
	} else {
		return ctx.Status(500).JSON(HttpError{Status: 500, Code: "internal-server", Message: err.Error()})
	}
}

type Messages map[string]string

type HttpError struct {
	Status  int         `binding:"required" json:"status"`
	Code    string      `binding:"required" json:"code"`
	Message interface{} `binding:"required" json:"message" oneOf:"string,Messages"`
}

func (e *HttpError) Error() string {
	if e.Message == nil {
		return ""
	}

	if str, ok := e.Message.(string); ok {
		return str
	}

	if strArr, ok := e.Message.([]string); ok {
		return strings.Join(strArr, ", ")
	}

	return "Unknown error"
}

func EntityNotFound(m interface{}) *HttpError {
	return &HttpError{Status: 404, Code: "entity-not-found", Message: m}
}

func BadRequest(m interface{}) *HttpError {
	return &HttpError{Status: 400, Code: "bad-request", Message: m}
}

func Unexpected(m interface{}) *HttpError {
	return &HttpError{Status: 500, Code: "internal-server", Message: m}
}

var (
	Unauthorized = &HttpError{Status: 401, Code: "unauthorized", Message: "You're not authorized"}
	Forbidden    = &HttpError{Status: 403, Code: "forbidden", Message: "You're not authorized"}
)
