package base

import (
	"errors"
	"fmt"
	"github.com/dotvezz/yoyo/internal/datatype"
)

// Errors returned from functions or methods in the `dialect` package MUST:
// 1. Wrap exactly one error from this list of vars
//    OR
// 2. Be an error from this list of vars
var (
	// InvalidDatatype is an error to use when a datatype referenced does not match anyything in the `datatype` package
	InvalidDatatype = errors.New("invalid datatype")
	// unsupportedDatatype is an error to use when a datatype is valid, but not supported in the current operation or
	// by the current dialect
	UnsupportedDatatype = errors.New("unsupported datatype")
)

func (*Base) InvalidDatatype(dt datatype.Datatype) error {
	return fmt.Errorf("%w: datatype `%s` is invalid", InvalidDatatype, dt)
}

func (d *Base) UnsupportedDatatype(dt datatype.Datatype) error {
	return fmt.Errorf("%w: %s does not support %s", UnsupportedDatatype, d.Dialect, dt)
}
