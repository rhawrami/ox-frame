package vector

import "github.com/rhawrami/ox-frame/ox/dtype"

// type StringVector represents a String Vector.
//
// String takes directly from Apache Arrow's string-type implementation.
type StringVector struct {
	DType     dtype.DataType
	validity  ValidityBitMap
	data      []byte  // all strings stored together; slice assumed to contain valid utf-8 sequences
	offsets   []int64 // offsets are length(len) + 1; e.g., final element takes two spots (start and end of last element)
	nullCount int
	len       int
}

func (v *StringVector) Type() dtype.DataType {
	return v.DType
}

func (v *StringVector) Len() int {
	return v.len
}

func (v *StringVector) NullCount() int {
	return v.nullCount
}

func (v *StringVector) Data() []byte {
	return v.data
}

func (v *StringVector) Offsets() []int64 {
	return v.offsets
}

func (v *StringVector) ValAt(i int) []byte {
	return v.data[v.offsets[i]:v.offsets[i+1]]
}

func (v *StringVector) Validity() ValidityBitMap {
	return v.validity
}

func (v *StringVector) IsNull(i int) bool {
	return v.validity.IsNull(i)
}

func (v *StringVector) IsNullBinary(i int) byte {
	return v.validity.IsNullBinary(i)
}

func (v *StringVector) DeepCopy() *StringVector {
	newData := make([]byte, len(v.data))
	newOffsets := make([]int64, len(v.offsets))
	newValidMapBuff := make([]byte, v.validity.Len())

	copy(newData, v.data)
	copy(newOffsets, v.offsets)
	copy(newValidMapBuff, v.validity.Buffer)

	return &StringVector{
		DType: v.DType,
		validity: ValidityBitMap{
			TrueLen:   v.len,
			NullCount: v.nullCount,
			Buffer:    newValidMapBuff,
		},
		data:      newData,
		offsets:   newOffsets,
		nullCount: v.nullCount,
		len:       v.len,
	}
}
