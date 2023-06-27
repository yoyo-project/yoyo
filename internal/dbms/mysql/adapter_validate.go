package mysql

import (
	"github.com/yoyo-project/yoyo/internal/datatype"
)

func (*adapter) SupportsDatatype(dt datatype.Datatype) bool {
	switch dt {
	case datatype.Integer,
		datatype.TinyInt,
		datatype.SmallInt,
		datatype.MediumInt,
		datatype.BigInt,
		datatype.Decimal,
		datatype.Real,
		datatype.Numeric,
		datatype.Double,
		datatype.Float,
		datatype.Varchar,
		datatype.Text,
		datatype.Blob,
		datatype.Enum,
		datatype.DateTime,
		datatype.Date,
		datatype.Timestamp,
		datatype.Binary,
		datatype.Year:
		return true
	}

	return false
}

func (*adapter) SupportsAutoIncrement() bool {
	return true
}
