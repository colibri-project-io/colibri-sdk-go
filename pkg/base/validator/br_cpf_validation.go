package validator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

const cpfRegex = `^\d{3}\.?\d{3}\.?\d{3}-?\d{2}$`

// brCPFValidation validates a CPF number.
func brCPFValidation(fl validator.FieldLevel) bool {
	r, _ := regexp.Compile(cpfRegex)

	return r.MatchString(fl.Field().String())
}
