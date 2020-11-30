package validation

import "github.com/dotvezz/yoyo/internal/datatype"

type Validator interface {
	SupportsDatatype(datatype datatype.Datatype) bool
}
