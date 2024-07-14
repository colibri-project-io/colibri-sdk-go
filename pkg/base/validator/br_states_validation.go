package validator

import (
	"slices"

	"github.com/go-playground/validator/v10"
)

// brStates is a list of Brazilian states.
var brStates = []string{
	"AC",
	"AL",
	"AP",
	"AM",
	"BA",
	"CE",
	"DF",
	"ES",
	"GO",
	"MA",
	"MT",
	"MS",
	"MG",
	"PA",
	"PB",
	"PR",
	"PE",
	"PI",
	"RJ",
	"RN",
	"RS",
	"RO",
	"RR",
	"SC",
	"SP",
	"SE",
	"TO",
}

// brStatesValidation validates a Brazilian state.
func brStatesValidation(fl validator.FieldLevel) bool {
	return slices.Contains(brStates, fl.Field().String())
}
