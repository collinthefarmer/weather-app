package validation

import (
	"errors"
)

var ErrValidation = errors.New("Could not validate.")

type Validates interface {
	Validate() (ValidationProblems, error)
}

type ValidationProblems = map[string]string
