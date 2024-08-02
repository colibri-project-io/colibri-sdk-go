package validator

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestBrStatesValidation(t *testing.T) {
	type testData struct {
		State    string `validate:"br_states"`
		Expected bool
	}

	tests := []testData{
		{State: "SP", Expected: true},
		{State: "RJ", Expected: true},
		{State: "MG", Expected: true},
		{State: "NY", Expected: false},
		{State: "CA", Expected: false},
	}

	validate := validator.New()
	validate.RegisterValidation("br_states", brStatesValidation)

	for _, test := range tests {
		err := validate.Var(test.State, "br_states")
		if test.Expected {
			assert.NoError(t, err)
		} else {
			assert.Error(t, err)
		}
	}
}
