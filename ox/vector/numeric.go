package vector

import (
	"github.com/rhawrami/ox-frame/ox/dtype"
)

// type Numeric includes all Go primitive numeric types
type Numeric interface {
	uint8 | uint16 | uint32 | uint64 |
		int8 | int16 | int32 | int64 | int |
		float32 | float64
}

// type NumericVector represents a Numeric Vector
type NumericVector[T Numeric] struct {
	DType     dtype.DataType
	validity  ValidityBitMap
	data      []T
	nullCount int
	len       int
}

func (v *NumericVector[T]) Type() dtype.DataType {
	return v.DType
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
		DType: v.DType,
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

// NewNumericVector returns a NumericVector
func NewNumericVector[T Numeric](data []T, validity []byte) *NumericVector[T] {
	var dt dtype.DataType
	switch any(data).(type) {
	case int, int64:
		dt = dtype.Int64{}
	default:
		dt = dtype.Float64{}
	}

	vm := ValidityBitMap{
		TrueLen:   len(data),
		NullCount: 0,
		Buffer:    validity,
	}
	return &NumericVector[T]{
		DType:     dt,
		validity:  vm,
		data:      data,
		nullCount: vm.CalcNullCount(),
		len:       len(data),
	}
}
