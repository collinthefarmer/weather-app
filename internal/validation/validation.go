package validation

type Validates interface {
	Validate() (ValidationProblems, error)
}

type ValidationProblems = map[string]string
