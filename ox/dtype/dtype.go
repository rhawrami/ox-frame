package dtype

// A DataType represents a column type supported by ox
type DataType interface {
	Type() LogicalType
	BitsReq() int
}

// LogicalType represents a logical type that is supported by ox
type LogicalType int

const (
	NULL LogicalType = iota

	BOOL

	UINT8

	INT8

	UINT16

	INT16

	UINT32

	INT32

	UINT64

	INT64

	FLOAT32

	FLOAT64

	STRING
)
