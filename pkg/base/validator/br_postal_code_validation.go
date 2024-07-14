package validator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

const brPostalCodeRegex = `^\d{5}[\-]?\d{3}$`

// brPostalCodeValidation validates a Brazilian postal code.
func brPostalCodeValidation(fl validator.FieldLevel) bool {
	r, _ := regexp.Compile(brPostalCodeRegex)

	return r.MatchString(fl.Field().String())
}
