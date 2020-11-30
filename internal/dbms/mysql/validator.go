package mysql

import "github.com/dotvezz/yoyo/internal/datatype"

type validator struct {
}

func (_ *validator) SupportsDatatype(dt datatype.Datatype) bool {
	switch dt {
	case datatype.Integer,
		datatype.TinyInt,
		datatype.SmallInt,
		datatype.MediumInt,
		datatype.BigInt,
		datatype.Decimal,
		datatype.Varchar,
		datatype.Text,
		datatype.Blob,
		datatype.Enum:
		return true
	}

	return false
}
