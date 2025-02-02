package errors

import (
	"errors"

	"gorm.io/gorm"
)

func EntityError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return EntityNotFound(err.Error())
	}

	return Unexpected(err.Error())
}
