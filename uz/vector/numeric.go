package vector

import (
	"github.com/rhawrami/uz-frame/uz/dtype"
)

// type Numeric includes all Go primitive numeric types
type Numeric interface {
	uint8 | uint16 | uint32 | uint64 |
		int8 | int16 | int32 | int64 | int |
		float32 | float64
}

// type NumericVector represents a Numeric Vector
type NumericVector[T Numeric] struct {
	dType     dtype.DataType
	validity  ValidityBitMap
	data      []T
	nullCount int
	len       int
}

func (v *NumericVector[T]) Type() dtype.DataType {
	return v.dType
}

func (v *NumericVector[T]) Len() int {
	return v.len
}

func (v *NumericVector[T]) NullCount() int {
	return v.nullCount
}

func (v *NumericVector[T]) Data() []T {
	return v.data
}

func (v *NumericVector[T]) ValAt(i int) T {
	return v.data[i]
}

func (v *NumericVector[T]) Validity() ValidityBitMap {
	return v.validity
}

func (v *NumericVector[T]) IsNull(i int) bool {
	return v.validity.IsNull(i)
}

func (v *NumericVector[T]) IsNullBinary(i int) byte {
	return v.validity.IsNullBinary(i)
}

func (v *NumericVector[T]) DeepCopy() *NumericVector[T] {
	newData := make([]T, v.len)
	newValidMapBuff := make([]byte, v.validity.Len())

	copy(newData, v.data)
	copy(newValidMapBuff, v.validity.Buffer)

	return &NumericVector[T]{
		dType: v.dType,
		validity: ValidityBitMap{
			TrueLen:   v.len,
			NullCount: v.nullCount,
			Buffer:    newValidMapBuff,
		},
		data:      newData,
		nullCount: v.nullCount,
		len:       v.len,
	}
}

// NumericVecFromComponents returns a NumericVector, given a Datatype, data, ValidityBitMap
func NumericVecFromComponents[T Numeric](dType dtype.DataType, data []T, validity ValidityBitMap) *NumericVector[T] {
	return &NumericVector[T]{
		dType:     dType,
		validity:  validity,
		data:      data,
		nullCount: validity.NullCount,
		len:       len(data),
	}
}

// NumericVecFromNums returns a NumericVector, given a slice of Numerics and a bool slice representing nulls
func NumericVecFromNums[T Numeric](data []T, validity []bool) *NumericVector[T] {
	validMap := ValidityBitMapFromBools(validity)
	dT := GetNumericDType(data[0])
	return &NumericVector[T]{
		dType:     dT,
		validity:  validMap,
		data:      data,
		nullCount: validMap.NullCount,
		len:       len(data),
	}
}

// GetNumericDType returns the Dtype of an element, given the element's native type
func GetNumericDType[T Numeric](x T) dtype.DataType {
	var dt dtype.DataType

	switch any(x).(type) {
	case int, int64:
		dt = dtype.Int64{}
	case int32:
		dt = dtype.Int32{}
	case int16:
		dt = dtype.Int16{}
	case int8:
		dt = dtype.Int8{}
	case float64:
		dt = dtype.Float64{}
	case float32:
		dt = dtype.Float32{}
	case uint64:
		dt = dtype.UInt64{}
	case uint32:
		dt = dtype.UInt32{}
	case uint16:
		dt = dtype.UInt16{}
	case uint8:
		dt = dtype.UInt8{}
	default:
		// non-implemented types
	}
	return dt
}
