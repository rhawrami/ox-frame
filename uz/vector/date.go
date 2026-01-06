package vector

import (
	"time"

	"github.com/rhawrami/uz-frame/uz/dtype"
)

const secsInOneDay int64 = 60 * 60 * 24

// type DateVector represents a Date Vector
type DateVector struct {
	dType     dtype.DataType
	validity  ValidityBitMap
	data      []int32
	nullCount int
	len       int
}

func (v *DateVector) Type() dtype.DataType {
	return v.dType
}

func (v *DateVector) Len() int {
	return v.len
}

func (v *DateVector) NullCount() int {
	return v.nullCount
}

func (v *DateVector) Data() []int32 {
	return v.data
}

func (v *DateVector) ValAt(i int) int32 {
	return v.data[i]
}

func (v *DateVector) Validity() ValidityBitMap {
	return v.validity
}

func (v *DateVector) IsNull(i int) bool {
	return v.validity.IsNull(i)
}

func (v *DateVector) IsNullBinary(i int) byte {
	return v.validity.IsNullBinary(i)
}

func (v *DateVector) DeepCopy() *DateVector {
	newData := make([]int32, v.len)
	newValidMapBuff := make([]byte, v.validity.Len())

	copy(newData, v.data)
	copy(newValidMapBuff, v.validity.Buffer)

	return &DateVector{
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

// DateVecFromComponents returns a DateVector, given data, and a ValidityBitMap
func DateVecFromComponents(data []int32, validity ValidityBitMap) *DateVector {
	return &DateVector{
		dType:     dtype.Date{},
		validity:  validity,
		data:      data,
		nullCount: validity.NullCount,
		len:       len(data),
	}
}

// DateVecFromStrings returns a DateVector, given date strings, a validity bool slice, and a string format
func DateVecFromStrings(data []string, validity []bool, layout string) *DateVector {
	validMap := ValidityBitMapFromBools(validity)
	dataI32 := make([]int32, len(data))
	for i, v := range data {
		var dI32 int32
		d, err := time.Parse(layout, v)

		if err != nil {
			validMap.SetNull(i)
			dI32 = 0
		} else {
			dI32 = int32(d.Unix() / secsInOneDay)
		}

		dataI32[i] = dI32
	}

	return &DateVector{
		dType:     dtype.Date{},
		validity:  validMap,
		data:      dataI32,
		nullCount: validMap.NullCount,
		len:       len(data),
	}
}
