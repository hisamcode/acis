package validator

import (
	"regexp"
	"slices"
)

// regular expression for checking email address
// this taken from https://html.spec.whatwg.org/#valid-e-mail-address
var (
	EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

type Validator struct {
	NonFieldErrors []string
	FieldErrors    map[string]string
}

// Creates a new validator instance with an empty errors map
func New() *Validator {
	return &Validator{FieldErrors: make(map[string]string)}
}

// return true if errors map doesn't contain any entries
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0
}

// Adds an error message to the map (so long as no entry already exists for the given key)
func (v *Validator) AddFieldError(key, message string) {
	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

// adds an error message to the map only if a validation check is not 'ok'
func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

// returns true if a spesific value is in a list of permitted values
func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}

// return true if a string value matches a spesific regexp pattern
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

// return true if string1 and string2 same otherwise false
func Equal(v1, v2 string) bool {
	return v1 == v2
}

// return true if all values in a slice are unique
func Unique[T comparable](values []T) bool {
	uniqueValues := make(map[T]bool)

	for _, value := range values {
		uniqueValues[value] = true
	}

	return len(values) == len(uniqueValues)
}
