package ftag

import "errors"

type withKind struct {
	underlying error
	tag        Kind
}

// Kind is a simple string to describe the category of an error.
type Kind string

// Implements all the interfaces for compatibility with the errors ecosystem.

func (e *withKind) Error() string  { return "Error Tag: " + string(e.tag) }
func (e *withKind) Cause() error   { return e.underlying }
func (e *withKind) Unwrap() error  { return e.underlying }
func (e *withKind) String() string { return e.Error() }

// Wrap wraps an error and gives it a distinct tag.
func Wrap(parent error, k Kind) error {
	if parent == nil {
		return nil
	}

	if k == "" {
		return parent
	}

	return &withKind{
		underlying: parent,
		tag:        k,
	}
}

// With implements the Fault Wrapper interface.
func With(k Kind) func(error) error {
	return func(err error) error {
		return Wrap(err, k)
	}
}

// Get extracts the error tag of an error chain. If there's no tag, returns nil.
func Get(err error) Kind {
	if err == nil {
		return None
	}

	for err != nil {
		if f, ok := err.(*withKind); ok {
			return f.tag
		}

		err = errors.Unwrap(err)
	}

	return Internal
}

func GetAll(err error) []Kind {
	if err == nil {
		return nil
	}

	var ks []Kind
	for err != nil {
		if f, ok := err.(*withKind); ok {
			ks = append(ks, f.tag)
		}

		err = errors.Unwrap(err)
	}

	return ks
}

// Is checks if the error has a specific kind.
func Is(err error, k Kind) bool {
	if err == nil {
		return false
	}

	for _, kind := range GetAll(err) {
		if kind == k {
			return true
		}
	}

	return false
}
