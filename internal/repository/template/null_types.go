package template

const NullTypeFile = `package nullable

import ("database/sql")

// First, re-export the types from database/sql

type Time = sql.NullTime
type Int16 = sql.NullInt16
type Int32 = sql.NullInt32
type Int64 = sql.NullInt64
type Bool = sql.NullBool
type Byte = sql.NullByte
type String = sql.NullString
type Float64 = sql.NullFloat64

// Later, export new types as needed...

`