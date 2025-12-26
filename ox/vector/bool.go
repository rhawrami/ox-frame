package vector

import (
	"github.com/rhawrami/ox-frame/ox/dtype"
)

// type BoolVector represents a Boolean Vector
type BoolVector struct {
	dType     dtype.DataType
	validity  ValidityBitMap
	data      []byte
	nullCount int
	len       int
}

func (v *BoolVector) Type() dtype.DataType {
	return v.dType
}

func (v *BoolVector) Len() int {
	return v.len
}

func (v *BoolVector) NullCount() int {
	return v.nullCount
}

func (v *BoolVector) Data() []byte {
	return v.data
}

func (v *BoolVector) ValAt(i int) bool {
	byteIdx, shiftBy := i/8, i%8
	return (v.data[byteIdx]>>byte(shiftBy))&byte(1) == 1
}

func (v *BoolVector) Validity() ValidityBitMap {
	return v.validity
}

func (v *BoolVector) IsNull(i int) bool {
	return v.validity.IsNull(i)
}

func (v *BoolVector) IsNullBinary(i int) byte {
	return v.validity.IsNullBinary(i)
}

func (v *BoolVector) DeepCopy() *BoolVector {
	newData := make([]byte, v.len)
	newValidMapBuff := make([]byte, v.validity.Len())

	copy(newData, v.data)
	copy(newValidMapBuff, v.validity.Buffer)

	return &BoolVector{
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

// BoolVecFromComponenets returns a BoolVector, given data, and a ValidityBitMap
func BoolVecFromComponenets(dType dtype.DataType, data []byte, validity ValidityBitMap) *BoolVector {
	return &BoolVector{
		dType:     dtype.Bool{},
		validity:  validity,
		data:      data,
		nullCount: validity.NullCount,
		len:       validity.TrueLen,
	}
}

// BoolVecFromBools returns a BoolVector, given a slice of bools and a bool slice representing nulls
func BoolVecFromBools(data []bool, validity []bool) *BoolVector {
	lenB := len(data)/8 + 1
	if len(data)%8 == 0 {
		lenB = lenB - 1
	}

	dataB := make([]byte, lenB)
	validityB := make([]byte, lenB)
	nNull := len(data)

	for i := 0; i < len(data); i++ {
		bIdx, shiftBy := i/8, i%8
		if data[i] {
			dataB[bIdx] = dataB[bIdx] | (1 << shiftBy)
		}
		if validity[i] {
			validityB[bIdx] = validityB[bIdx] | (1 << shiftBy)
			nNull -= 1
		}
	}

	validMap := ValidityBitMap{
		TrueLen:   len(data),
		NullCount: nNull,
		Buffer:    validityB,
	}

	return &BoolVector{
		dType:     dtype.Bool{},
		validity:  validMap,
		data:      dataB,
		nullCount: nNull,
		len:       len(data),
	}
}
