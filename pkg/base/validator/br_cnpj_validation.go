package validator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

const cnpjRegex = `^\d{2}\.?\d{3}\.?\d{3}/?\d{4}-?\d{2}$`

// brCNPJValidation validates a CNPJ number.
func brCNPJValidation(fl validator.FieldLevel) bool {
	r, _ := regexp.Compile(cnpjRegex)

	return r.MatchString(fl.Field().String())
}
