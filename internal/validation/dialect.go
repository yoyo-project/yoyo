package validation

import "github.com/dotvezz/yoyo/internal/datatype"

// Validator provides an interface for DBMS-specific validations
type Validator interface {
	// SupportsDatatype takes a datatype and returns true if the underlying DBMS supports the datatype
	SupportsDatatype(datatype datatype.Datatype) bool
}
