package validator

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestBrCNPJValidation(t *testing.T) {
	type testData struct {
		CNPJ     string `validate:"br_cnpj"`
		Expected bool
	}

	tests := []testData{
		{CNPJ: "12.345.678/0001-90", Expected: true},
		{CNPJ: "12.345.678/0001-9a", Expected: false},
		{CNPJ: "12.345.678/0001-900", Expected: false},
		{CNPJ: "12.345.678/0001-9000", Expected: false},
		{CNPJ: "12.345.678/0001-90000", Expected: false},
		{CNPJ: "12345678000190", Expected: true},
		{CNPJ: "1234567800019a", Expected: false},
		{CNPJ: "123456780001900", Expected: false},
		{CNPJ: "1234567800019000", Expected: false},
		{CNPJ: "12345678000190000", Expected: false},
	}

	validate := validator.New()
	validate.RegisterValidation("br_cnpj", brCNPJValidation)

	for _, test := range tests {
		err := validate.Var(test.CNPJ, "br_cnpj")
		if test.Expected {
			assert.NoError(t, err)
		} else {
			assert.Error(t, err)
		}
	}
}
