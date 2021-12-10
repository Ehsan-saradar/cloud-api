package validation

import (
	"api.cloud.io/pkg/errors"
	"github.com/asaskevich/govalidator"
)

// ValidateStruct validate a struct by govalidator tags.
func ValidateStruct(s interface{}) error {
	isValid, err := govalidator.ValidateStruct(s)
	if !isValid || s == nil {
		return errors.ErrInvalidParams(err, govalidator.ErrorsByField(err))
	}

	return nil
}
