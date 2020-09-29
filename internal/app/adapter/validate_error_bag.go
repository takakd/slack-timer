package adapter

import "proteinreminder/internal/pkg/collection"

// Validate Error type
type ErrorType string

// Defined Error Type
const (
	Empty ErrorType = ""
)

// Set of validate error
type ValidateErrorBag struct {
	errors map[string]*ValidateError
}

func NewValidateErrorBag() *ValidateErrorBag {
	b := &ValidateErrorBag{}
	b.errors = make(map[string]*ValidateError)
	return b
}

type ValidateError struct {
	Name    string
	Summary string
	TypeSet *collection.Set
}

func (b *ValidateErrorBag) SetError(name, summary string, errorType ErrorType) {
	var error *ValidateError
	error, errorExists := b.errors[name]
	if !errorExists {
		error = &ValidateError{
			Name:    name,
			Summary: summary,
		}
		error.TypeSet = collection.NewSet()
	}

	error.TypeSet.Set(errorType)

	if summary != "" {
		error.Summary = summary
	}

	b.errors[name] = error
}

// Check if error exists.
// true: exist, false: not exist.
func (b *ValidateErrorBag) ContainsError(name string, errorType ErrorType) bool {
	error, errorExists := b.errors[name]
	if !errorExists {
		return false
	}

	return error.TypeSet.Contains(errorType)
}

func (b *ValidateErrorBag) GetError(name string) (*ValidateError, bool) {
	error, errorExists := b.errors[name]
	return error, errorExists
}

func (b *ValidateErrorBag) GetErrors() []*ValidateError{
	errors := make([]*ValidateError, 0, len(b.errors))
	for _, v := range b.errors {
		errors = append(errors, v)
	}
	return errors
}
