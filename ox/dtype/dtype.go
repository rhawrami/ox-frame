package dtype

// Note: The code below is taken from the Apache Arrow-Go implementation:
// https://github.com/apache/arrow-go/blob/main/arrow/datatype.go
// Only primitive types, and strings, are currently supported by Ox

// A DataType represents a column type supported by ox
type DataType interface {
	Type() LogicalType
	BitsReq() int
}

// LogicalType represents a logical type that is supported by ox
type LogicalType int

const (
	// NULL type having no physical storage
	NULL LogicalType = iota

	// BOOL is a 1 bit, LSB bit-packed ordering
	BOOL

	// UINT8 is an Unsigned 8-bit little-endian integer
	UINT8

	// INT8 is a Signed 8-bit little-endian integer
	INT8

	// UINT16 is an Unsigned 16-bit little-endian integer
	UINT16

	// INT16 is a Signed 16-bit little-endian integer
	INT16

	// UINT32 is an Unsigned 32-bit little-endian integer
	UINT32

	// INT32 is a Signed 32-bit little-endian integer
	INT32

	// UINT64 is an Unsigned 64-bit little-endian integer
	UINT64

	// INT64 is a Signed 64-bit little-endian integer
	INT64

	// FLOAT32 is a 4-byte floating point value
	FLOAT32

	// FLOAT64 is an 8-byte floating point value
	FLOAT64

	// STRING is a UTF8 variable-length string
	STRING

	// String (UTF8) view type with 4-byte prefix and inline
	// small string optimizations
	STRING_VIEW

	// DATE32 is int32 days since the UNIX epoch
	DATE32

	// DATE64 is int64 milliseconds since the UNIX epoch
	DATE64
)

// IsPrimitive returns true if the provided type ID represents a fixed width
// primitive type.
func IsPrimitive(t LogicalType) bool {
	switch t {
	case BOOL,
		UINT8, INT8, UINT16, INT16, UINT32, INT32, UINT64, INT64,
		FLOAT32, FLOAT64,
		DATE32, DATE64:
		return true
	}
	return false
}

// IsUnsignedInteger is a helper that returns true if the type ID provided is
// one of the uint integral types (uint8, uint16, uint32, uint64)
func IsUnsignedInteger(t LogicalType) bool {
	switch t {
	case UINT8, UINT16, UINT32, UINT64:
		return true
	}
	return false
}

// IsSignedInteger is a helper that returns true if the type ID provided is
// one of the int integral types (int8, int16, int32, int64)
func IsSignedInteger(t LogicalType) bool {
	switch t {
	case INT8, INT16, INT32, INT64:
		return true
	}
	return false
}

// IsFloating is a helper that returns true if the type ID provided is
// one of Float16, Float32, or Float64
func IsFloating(t LogicalType) bool {
	switch t {
	case FLOAT32, FLOAT64:
		return true
	}
	return false
}
