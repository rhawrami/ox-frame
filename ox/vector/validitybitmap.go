package vector

import "math/bits"

// ValidityBitMap represents a bitmap of null values corresponding to
// a dataframe column.
//
// Taken after Apache Arrow's validity bitmap
type ValidityBitMap struct {
	TrueLen   int    // actual number of elements represented by bitmap
	NullCount int    // number of null elements
	Buffer    []byte // null bitmap, little-endian; 1 == NOT NULL
}

// Len returns the byte-length (e.g., not "true" length) of the ValidityBitMap
func (m ValidityBitMap) Len() int {
	return len(m.Buffer)
}

// LookUp returns whether a corresponding column record `i` is null.
//
// Assumes record at `i` exists
func (m ValidityBitMap) IsNull(i int) bool {
	byteIdx, shiftBy := i/8, i%8
	return (m.Buffer[byteIdx]>>byte(shiftBy))&byte(1) == 0
}

// IsNullBinary returns 1 if an element is null, and 0 otherwise
func (m ValidityBitMap) IsNullBinary(i int) byte {
	byteIdx, shiftBy := i/8, i%8
	return (m.Buffer[byteIdx]>>byte(shiftBy))&byte(1) ^ byte(1)
}

// SetNull sets a corresponding column record to null
func (m ValidityBitMap) SetNull(i int) {
	byteIdx, shiftBy := i/8, i%8
	m.Buffer[byteIdx] = m.Buffer[byteIdx] &^ (1 << shiftBy)
}

// SetNotNull sets a corresponding column record to not-null
func (m ValidityBitMap) SetNotNull(i int) {
	byteIdx, shiftBy := i/8, i%8
	m.Buffer[byteIdx] = m.Buffer[byteIdx] | (1 << shiftBy)
}

// CalcNullCount manually calculates the null count of a ValidityBitMap
// rather than relying on the internal null count field
func (m ValidityBitMap) CalcNullCount() int {
	n := m.TrueLen
	for i := 0; i < len(m.Buffer); i++ {
		n -= bits.OnesCount8(m.Buffer[i])
	}
	return n
}

// ValidityBitMapFromBools returns a new ValidityBitMap, taking in a
// boolean slice as input.
func ValidityBitMapFromBools(b []bool) ValidityBitMap {
	// more likely than not to not be div by 8
	// add remainder first, subtract if needed
	lenMap := len(b)/8 + 1
	if len(b)%8 == 0 {
		lenMap = lenMap - 1
	}

	buff := make([]byte, lenMap)
	nNull := len(b)

	for i := 0; i < len(b); i++ {
		if b[i] {
			bIdx, shiftBy := i/8, i%8
			// flip zero bit
			buff[bIdx] = buff[bIdx] | (1 << shiftBy)
			nNull -= 1
		}
	}

	return ValidityBitMap{
		TrueLen:   len(b),
		NullCount: nNull,
		Buffer:    buff,
	}
}
