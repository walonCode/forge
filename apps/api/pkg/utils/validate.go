package utils

import "github.com/go-playground/validator/v10"

// validate is a process-wide validator instance (it is safe for concurrent use
// and caches struct metadata, so it should be reused rather than recreated).
var validate = validator.New()

// ValidateStruct validates a struct against its `validate` tags and returns an
// error describing the first failing field, or nil when the struct is valid.
func ValidateStruct(s any) error {
	return validate.Struct(s)
}
