package vector

import "github.com/rhawrami/rok-frame/rok/dtype"

// type StringVector represents a String Vector.
//
// String takes directly from Apache Arrow's string-type implementation.
type StringVector struct {
	dType     dtype.DataType
	validity  ValidityBitMap
	data      []byte  // all strings stored together; slice assumed to contain valid utf-8 sequences
	offsets   []int64 // offsets are length(len) + 1; e.g., final element takes two spots (start and end of last element)
	nullCount int
	len       int
}

func (v *StringVector) Type() dtype.DataType {
	return v.dType
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

func (v *StringVector) StringValAt(i int) string {
	return string(v.data[v.offsets[i]:v.offsets[i+1]])
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
		dType: v.dType,
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

// StringVecFromComponents returns a StringVector, given data, offsets and a ValidityBitMap
func StringVecFromComponents(data []byte, offsets []int64, validity ValidityBitMap) *StringVector {
	return &StringVector{
		dType:     dtype.String{},
		validity:  validity,
		data:      data,
		offsets:   offsets,
		nullCount: validity.NullCount,
		len:       len(offsets) - 1, // last offset element is starting place of imaginary N+1'th element
	}
}

// StringVecFromStrings returns a StringVector, given a slice of strings and a slice of bools representing nulls
func StringVecFromStrings(data []string, validity []bool) *StringVector {
	validMap := ValidityBitMapFromBools(validity)
	// assume strings are 4 bytes on average; reduce reallocs
	dataInBytes := make([]byte, 0, len(data)*4)
	offsets := make([]int64, len(data)+1)
	for i, v := range data {
		offsets[i] = int64(len(dataInBytes))
		dataInBytes = append(dataInBytes, []byte(v)...)
	}
	// final offset element
	offsets[len(offsets)-1] = int64(len(dataInBytes))

	return &StringVector{
		dType:     dtype.String{},
		validity:  validMap,
		data:      dataInBytes,
		offsets:   offsets,
		nullCount: validMap.NullCount,
		len:       len(offsets) - 1,
	}
}
