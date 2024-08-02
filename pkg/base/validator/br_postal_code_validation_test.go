package validator

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestBrPostalCodeValidation(t *testing.T) {
	type testData struct {
		PostalCode string
		Expected   bool
	}

	tests := []testData{
		{PostalCode: "12345-678", Expected: true},
		{PostalCode: "12345678", Expected: true},
		{PostalCode: "12345", Expected: false},
		{PostalCode: "123456789", Expected: false},
		{PostalCode: "ABCDE-FGH", Expected: false},
		{PostalCode: "ABCDEFGH", Expected: false},
	}

	validate := validator.New()
	validate.RegisterValidation("br_postal_code", brPostalCodeValidation)

	for _, test := range tests {
		err := validate.Var(test.PostalCode, "br_postal_code")
		if test.Expected {
			assert.NoError(t, err)
		} else {
			assert.Error(t, err)
		}
	}
}
