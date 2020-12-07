// Package validator provides features that validate request parameters.
package validator

import (
	"slacktimer/internal/pkg/collection"
)

// ValidateErrorBag represents validate error
type ValidateErrorBag struct {
	errors map[string]*ValidateError
}

// NewValidateErrorBag create new struct.
func NewValidateErrorBag() *ValidateErrorBag {
	b := &ValidateErrorBag{}
	b.errors = make(map[string]*ValidateError)
	return b
}

// ValidateError represents validation errors of one value.
type ValidateError struct {
	Name    string
	Summary string
	TypeSet *collection.Set
}

// SetError set the validation error to value identified by name.
func (b *ValidateErrorBag) SetError(name, summary string, err error) {
	var error *ValidateError
	error, errorExists := b.errors[name]
	if !errorExists {
		error = &ValidateError{
			Name:    name,
			Summary: summary,
		}
		error.TypeSet = collection.NewSet()
	}

	error.TypeSet.Set(err)

	if summary != "" {
		error.Summary = summary
	}

	b.errors[name] = error
}

// ContainsError checks if error exists.
// true: exist, false: not exist.
func (b ValidateErrorBag) ContainsError(name string, err error) bool {
	error, errorExists := b.errors[name]
	if !errorExists {
		return false
	}

	return error.TypeSet.Contains(err)
}

// GetError returns validation error. if errors does not exists, returns false.
func (b ValidateErrorBag) GetError(name string) (*ValidateError, bool) {
	error, errorExists := b.errors[name]
	return error, errorExists
}

// GetErrors returns all validation errors.
func (b ValidateErrorBag) GetErrors() []*ValidateError {
	errors := make([]*ValidateError, 0, len(b.errors))
	for _, v := range b.errors {
		errors = append(errors, v)
	}
	return errors
}
