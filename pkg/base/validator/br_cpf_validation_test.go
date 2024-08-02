package validator

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestBrCPFValidation(t *testing.T) {
	type testData struct {
		CPF      string `validate:"br_cpf"`
		Expected bool
	}

	tests := []testData{
		{CPF: "123.456.789-00", Expected: true},
		{CPF: "111.222.333-44", Expected: true},
		{CPF: "999.888.777-66", Expected: true},
		{CPF: "12345678900", Expected: true},
		{CPF: "1234567890", Expected: false},
		{CPF: "123456789000", Expected: false},
		{CPF: "abc.def.ghi-jk", Expected: false},
		{CPF: "abcdefghijk", Expected: false},
	}

	validate := validator.New()
	validate.RegisterValidation("br_cpf", brCPFValidation)

	for _, test := range tests {
		err := validate.Var(test.CPF, "br_cpf")
		if test.Expected {
			assert.NoError(t, err)
		} else {
			assert.Error(t, err)
		}
	}
}
