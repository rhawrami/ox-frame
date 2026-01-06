package dtype

// Bool represents a boolean
type Bool struct{}

func (x Bool) Type() LogicalType { return BOOL }
func (x Bool) String() string    { return "bool" }
func (x Bool) BitsReq() int      { return 1 }

// UInt8 represents an 8-bit unsigned integer
type UInt8 struct{}

func (x UInt8) Type() LogicalType { return UINT8 }
func (x UInt8) String() string    { return "uint8" }
func (x UInt8) BitsReq() int      { return 8 }

// UInt16 represents an 16-bit unsigned integer
type UInt16 struct{}

func (x UInt16) Type() LogicalType { return UINT16 }
func (x UInt16) String() string    { return "uint16" }
func (x UInt16) BitsReq() int      { return 16 }

// UInt32 represents an 32-bit unsigned integer
type UInt32 struct{}

func (x UInt32) Type() LogicalType { return UINT32 }
func (x UInt32) String() string    { return "uint32" }
func (x UInt32) BitsReq() int      { return 32 }

// UInt64 represents an 64-bit unsigned integer
type UInt64 struct{}

func (x UInt64) Type() LogicalType { return UINT64 }
func (x UInt64) String() string    { return "uint64" }
func (x UInt64) BitsReq() int      { return 64 }

// Int8 represents an 8-bit unsigned integer
type Int8 struct{}

func (x Int8) Type() LogicalType { return INT8 }
func (x Int8) String() string    { return "int8" }
func (x Int8) BitsReq() int      { return 8 }

// Int16 represents an 16-bit unsigned integer
type Int16 struct{}

func (x Int16) Type() LogicalType { return INT16 }
func (x Int16) String() string    { return "int16" }
func (x Int16) BitsReq() int      { return 16 }

// Int32 represents an 32-bit unsigned integer
type Int32 struct{}

func (x Int32) Type() LogicalType { return INT32 }
func (x Int32) String() string    { return "int32" }
func (x Int32) BitsReq() int      { return 32 }

// Int64 represents an 64-bit unsigned integer
type Int64 struct{}

func (x Int64) Type() LogicalType { return INT64 }
func (x Int64) String() string    { return "int64" }
func (x Int64) BitsReq() int      { return 64 }

// Float32 represents a 32-bit floating point
type Float32 struct{}

func (x Float32) Type() LogicalType { return FLOAT32 }
func (x Float32) String() string    { return "float32" }
func (x Float32) BitsReq() int      { return 32 }

// Float64 represents a 64-bit floating point
type Float64 struct{}

func (x Float64) Type() LogicalType { return FLOAT64 }
func (x Float64) String() string    { return "float64" }
func (x Float64) BitsReq() int      { return 64 }

// String represents a string
type String struct{}

func (x String) Type() LogicalType { return STRING }
func (x String) String() string    { return "string" }
func (x String) BitsReq() int      { return -1 } // come back to this

// Date represents a date
type Date struct{}

func (x Date) Type() LogicalType { return DATE }
func (x Date) String() string    { return "date" }
func (x Date) BitsReq() int      { return 32 } // come back to this
