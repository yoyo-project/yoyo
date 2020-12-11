package postgres

import "github.com/dotvezz/yoyo/internal/datatype"

type validator struct {
}

func (*validator) SupportsDatatype(dt datatype.Datatype) bool {
	switch dt {
	case datatype.Integer,
		datatype.SmallInt,
		datatype.MediumInt,
		datatype.BigInt,
		datatype.Decimal,
		datatype.Varchar,
		datatype.Text,
		datatype.TinyText,
		datatype.MediumText,
		datatype.LongText,
		datatype.Char,
		datatype.Blob,
		datatype.Enum,
		datatype.Boolean:
		return true
	}

	return false
}
