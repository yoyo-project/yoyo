package datatype

import (
	"gopkg.in/yaml.v3"
	"strings"
)

// Datatype is used to encode information about types for use in repository or validation
// The least-significant 8 bits are reserved for general metadata
// The next 16 bits are not currently used. They were historically reserved for DBMS support in the early concept stage.
// The next 8 bits are reserved for unique type identification
// The last 32 bits are not currently used
type Datatype uint64

// These are the actual Datatype constants with all the metadata and unique identifiers encoded into them
const (
	Integer    = idInteger | metaNumeric | metaInteger | metaSignable | metaHasGoUnisgned
	TinyInt    = idTinyInt | metaNumeric | metaInteger | metaSignable | metaHasGoUnisgned
	SmallInt   = idSmallInt | metaNumeric | metaInteger | metaSignable | metaHasGoUnisgned
	MediumInt  = idMediumInt | metaNumeric | metaInteger | metaSignable | metaHasGoUnisgned
	BigInt     = idBigInt | metaNumeric | metaInteger | metaSignable | metaHasGoUnisgned
	Decimal    = idDecimal | metaNumeric | metaSignable | metaRequiresParams
	Varchar    = idVarchar | metaString
	Text       = idText | metaString
	TinyText   = idTinyText | metaString
	MediumText = idMediumText | metaString
	LongText   = idLongText | metaString
	Char       = idChar | metaString
	Blob       = idBlob | metaBinary
	Enum       = idEnum | metaString | metaRequiresParams
	Boolean    = idBoolean
	Date       = idDate | metaTime
	DateTime   = idDateTime | metaTime
	Timestamp  = idTimestamp | metaTime
)

// These are the string representations of datatypes
const (
	integer    = "INTEGER" // yoyo considers "INTEGER" to be the canonical string, however
	sint       = "INT"     // it still accepts "INT" as an alias and canonicalizes it to "INTEGER"
	tinyint    = "TINYINT"
	smallint   = "SMALLINT"
	mediumint  = "MEDIUMINT"
	bigint     = "BIGINT"
	decimal    = "DECIMAL"
	varchar    = "VARCHAR"
	text       = "TEXT"
	tinytext   = "TINYTEXT"
	mediumtext = "MEDIUMTEXT"
	longtext   = "LONGTEXT"
	char       = "CHAR"
	blob       = "BLOB"
	enum       = "ENUM"
	boolean    = "BOOLEAN" // yoyo considers "BOOLEAN" to be the canonical string, however
	sbool      = "BOOL"    // it still accepts "BOOL" as an alias and canonicalizes it to "BOOLEAN"

	goInt64   = "int64"
	goInt32   = "int32"
	goInt16   = "int16"
	goInt8    = "int8"
	goFloat64 = "float64"
	goString  = "string"
	goBool    = "bool"
	goRune    = "rune"
	goBlob    = "[]byte"
)

// UnmarshalYAML provides an implementation for yaml/v2.Unmarshaler to parse the yaml config
func (dt *Datatype) UnmarshalYAML(value *yaml.Node) (err error) {
	*dt, err = FromString(value.Value)
	return err
}

// MarshalYAML provides an implementation for yaml/v2.Marshaler, returns a string representation of the Datatype
func (dt Datatype) MarshalYAML() (interface{}, error) {
	return strings.ToLower(dt.String()), nil
}

func (dt Datatype) String() (s string) {
	switch dt {
	case Integer:
		s = integer
	case TinyInt:
		s = tinyint
	case SmallInt:
		s = smallint
	case MediumInt:
		s = mediumint
	case BigInt:
		s = bigint
	case Decimal:
		s = decimal
	case Varchar:
		s = varchar
	case Text:
		s = text
	case TinyText:
		s = tinytext
	case MediumText:
		s = mediumtext
	case LongText:
		s = longtext
	case Char:
		s = char
	case Blob:
		s = blob
	case Enum:
		s = enum
	case Boolean:
		s = boolean
	default:
		s = "NONE"
	}
	return s
}

func (dt Datatype) GoTypeString() (s string) {
	switch dt {
	case Integer:
		s = goInt32
	case TinyInt:
		s = goInt8
	case SmallInt:
		s = goInt16
	case MediumInt:
		s = goInt32
	case BigInt:
		s = goInt64
	case Decimal:
		s = goFloat64
	case Varchar:
		s = goString
	case Text:
		s = goString
	case TinyText:
		s = goString
	case MediumText:
		s = goString
	case LongText:
		s = goString
	case Char:
		s = goRune
	case Blob:
		s = goBlob
	case Enum:
		s = goString
	case Boolean:
		s = goBool
	default:
		s = "NONE"
	}
	return s
}

// IsInt returns true if the Datatype is an integer type
func (dt Datatype) IsInt() bool {
	return dt&metaInteger > 0
}

// IsNumeric returns true if the Datatype is a numeric type.
func (dt Datatype) IsNumeric() bool {
	return dt&metaNumeric > 0
}

// IsBinary returns true if the Datatype is a blob/binary type.
func (dt Datatype) IsBinary() bool {
	return dt&metaBinary > 0
}

// RequiresParams returns true if the Datatype requires parameters in SQL syntax.
// The `(8, 5)` in SQL's `DECIMAL(8, 5)` for example
func (dt Datatype) RequiresParams() bool {
	return dt&metaRequiresParams > 0
}

// IsString returns true if the Datatype is a string type
func (dt Datatype) IsString() bool {
	return dt&metaString > 0
}

// IsSignable returns true if the Datatype can be stored as either a signed or unsigned value
func (dt Datatype) IsSignable() bool {
	return dt&metaSignable > 0
}

// HasGoUnsigned returns true if the Datatype has an unsigned variant Go type like int and uint
func (dt Datatype) HasGoUnsigned() bool {
	return dt&metaHasGoUnisgned > 0
}

// IsTime returns true if the Datatype can be stored as either a signed or unsigned value
func (dt Datatype) IsTime() bool {
	return dt&metaTime > 0
}

// FromString returns the decoded Datatype, and an error if the in string is invalid or unknown
func FromString(in string) (dt Datatype, err error) {
	switch strings.ToUpper(strings.Split(in, "(")[0]) {
	case integer, sint:
		dt = Integer
	case bigint:
		dt = BigInt
	case mediumint:
		dt = MediumInt
	case smallint:
		dt = SmallInt
	case tinyint:
		dt = TinyInt
	case decimal:
		dt = Decimal
	case varchar:
		dt = Varchar
	case text:
		dt = Text
	case tinytext:
		dt = TinyText
	case mediumtext:
		dt = MediumText
	case longtext:
		dt = LongText
	case char:
		dt = Char
	case blob:
		dt = Blob
	case enum:
		dt = Enum
	case boolean, sbool:
		dt = Boolean
	default:
		err = ErrUnknownDatatype
	}

	return dt, err
}

// These metadata are general metadata to describe the data type
// 8 bits are reserved for this
const (
	metaNumeric Datatype = 1 << iota
	metaInteger
	metaBinary
	metaString
	metaTime
	metaSignable
	metaHasGoUnisgned
	metaRequiresParams
)

// These are the unique type identifiers
// Unlike the others, these are not single-bit flags
const (
	idInteger Datatype = (iota + 1) << 24
	idTinyInt
	idSmallInt
	idMediumInt
	idBigInt
	idDecimal
	idVarchar
	idText
	idTinyText
	idMediumText
	idLongText
	idChar
	idBlob
	idEnum
	idBoolean
	idDate
	idDateTime
	idTimestamp
)
