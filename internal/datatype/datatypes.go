package datatype

import "strings"

// Datatype is used to encode information about types for use in generation or validation
// The least-significant 8 bits are reserved for general metadata
// The next 16 bits are not currently used. They were historically reserved for DBMS support in the early concept stage.
// The next 8 bits are reserved for unique type identification
// The last 32 bits are not currently used
type Datatype uint64

// These are the actual Datatype constants with all the metadata and unique identifiers encoded into them
const (
	Integer    = idInteger | metaNumeric | metaInteger | metaSignable
	TinyInt    = idTinyInt | metaNumeric | metaInteger | metaSignable
	SmallInt   = idSmallInt | metaNumeric | metaInteger | metaSignable
	MediumInt  = idMediumInt | metaNumeric | metaInteger | metaSignable
	BigInt     = idBigInt | metaNumeric | metaInteger | metaSignable
	Decimal    = idDecimal | metaNumeric | metaSignable | metaRequiresScale
	Varchar    = idVarchar | metaString
	Text       = idText | metaString
	TinyText   = idTinyText | metaString
	MediumText = idMediumText | metaString
	LongText   = idLongText | metaString
	Char       = idChar | metaString
	Blob       = idBlob | metaBinary
	Enum       = idEnum | metaString | metaRequiresScale
	Boolean    = idBoolean
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
)

// UnmarshalYAML provides an implementation for yaml/v2.Unmarshaler to parse the yaml config
func (dt *Datatype) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	err := unmarshal(&s)
	if err != nil {
		return err
	}

	*dt, err = FromString(s)
	return err
}

func (dt Datatype) MarshalYAML() (interface{}, error) {
	return strings.ToLower(dt.String()), nil
}

func (dt Datatype) String() string {
	switch dt {
	case Integer:
		return integer
	case TinyInt:
		return tinyint
	case SmallInt:
		return smallint
	case MediumInt:
		return mediumint
	case BigInt:
		return bigint
	case Decimal:
		return decimal
	case Varchar:
		return varchar
	case Text:
		return text
	case TinyText:
		return tinytext
	case MediumText:
		return mediumtext
	case LongText:
		return longtext
	case Char:
		return char
	case Blob:
		return blob
	case Enum:
		return enum
	case Boolean:
		return boolean
	default:
		return "NONE"
	}
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

// RequiresScale returns true if the Datatype requires a range.
// The `(8, 5)` in MySQL's `DECIMAL(8, 5)` is a range, as far as yoyo is concerned
func (dt Datatype) RequiresScale() bool {
	return dt&metaRequiresScale > 0
}

// IsString returns true if the Datatype is a string type
func (dt Datatype) IsString() bool {
	return dt&metaString > 0
}

// IsSignable returns true if the Datatype can be stored as either a signed or unsigned value
func (dt Datatype) IsSignable() bool {
	return dt&metaSignable > 0
}

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
		err = UnknownDatatype
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
	metaSignable
	metaRequiresScale
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
)
