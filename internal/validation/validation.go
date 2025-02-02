package validation

import (
	"reflect"
	"strings"

	"github.com/guemidiborhane/factorydigitale.tech/internal/errors"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type (
	ErrorResponse struct {
		Error       bool
		FailedField string
		Tag         string
		Value       interface{}
	}

	CustomValidator struct {
		validator *validator.Validate
	}
)

var validate = validator.New(validator.WithRequiredStructEnabled())

func (v CustomValidator) Validate(data interface{}) []ErrorResponse {
	errors := []ErrorResponse{}

	errs := validate.Struct(data)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			// In this case data object is actually holding the User struct
			var elem ErrorResponse

			elem.FailedField = err.Field() // Export struct field name
			elem.Tag = err.Tag()           // Export struct tag
			elem.Value = err.Value()       // Export field value
			elem.Error = true

			errors = append(errors, elem)
		}
	}

	return errors
}

func Validate(body interface{}) error {
	if err := Validation.Validate(body); len(err) > 0 && err[0].Error {
		errMsgs := fiber.Map{}

		for _, err := range err {
			errMsgs[err.FailedField] = err.Tag
		}

		return errors.BadRequest(errMsgs)
	}
	return nil
}

var Validation *CustomValidator

func Setup() {
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	Validation = &CustomValidator{
		validator: validate,
	}
}
