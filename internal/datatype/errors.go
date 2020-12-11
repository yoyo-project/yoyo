package datatype

import "errors"

// ErrUnknownDatatype is the error returned by functions when an unknown or invalid datatype is used
var ErrUnknownDatatype = errors.New("invalid datatype")
