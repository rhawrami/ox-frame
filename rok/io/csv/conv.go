package csv

import (
	"time"

	"github.com/rhawrami/rok-frame/rok/vector"
)

type parsedRes interface {
	null() bool
}

// numericRes returns a parsed Numeric value, given a byte slice input
//
// isNull == true in the following cases:
// - empty slice (e.g., len(b) == 0)
// - non-numeric byte (excluding decimal points, and positive-negative signs)
type numericRes[T vector.Numeric] struct {
	val    T
	isNull bool
}

func (r numericRes[T]) null() bool { return r.isNull }

// strRes returns a parsed string value, given a byte slice input
//
// isNull == true when the slice is empty (e.g., len(b) == 0)
type strRes struct {
	val    []byte
	isNull bool
}

func (r strRes) null() bool { return r.isNull }

// dateRes returns a parsed date value, given a byte slice input
//
// isNull == true in the following cases:
// - empty slice (e.g., len(b) == 0)
// - byte slice is not able to be parsed into a time.Time value
// -	this can be triggered when the slice has ordinal suffixes (e.g., 2nd, 1st)
// -	or periods following an abbreviated month (e.g. Jan.)
type dateRes struct {
	val    int32
	isNull bool
}

func (r dateRes) null() bool { return r.isNull }

// boolRes returns a parsed boolean value, given a byte slice input
//
// isNull == true when the input is not one of the follownig formats:
// [`t`, `f`, `T`,`t`, `true`, `false`, `True`, `False`]
type boolRes struct {
	val    bool
	isNull bool
}

func (r boolRes) null() bool { return r.isNull }

// bToInt64 converts a byte slice to a 64-bit integer
// e.g. []byte("-4820") => int64(4820)
func bToInt64(b []byte) parsedRes {
	return bToSignedInteger[int64](b)
}

// bToInt32 converts a byte slice to a 32-bit integer
// e.g. []byte("-4820") => int32(4820)
func bToInt32(b []byte) parsedRes {
	return bToSignedInteger[int32](b)
}

// bToFloat64 converts a byte slice to a 64-bit floating-point
// e.g. []byte("+4820.7893") => float64(4820.7893)
func bToFloat64(b []byte) parsedRes {
	return bToFloatingPoint[float64](b)
}

// bToFloat32 converts a byte slice to a 32-bit floating-point
// e.g. []byte("+4820.7893") => float32(4820.7893)
func bToFloat32(b []byte) parsedRes {
	return bToFloatingPoint[float32](b)
}

// bToStr "converts" a byte slice to a string-type
// note: there isn't any actual conversion, as StringVector stores
// data in a contiguous byte slice
func bToStr(b []byte) parsedRes {
	var res strRes = strRes{val: []byte(""), isNull: true}
	if len(b) == 0 {
		return res
	}
	res.val, res.isNull = b, false
	return res
}

// bToNYearMonthDay converts a byte slice to a date type (N days since Unix epoch, stored as int32)
// Follows a YYYY-MM-DD format; e.g., "2006-01-02"
func bToNYearMonthDay(b []byte) parsedRes {
	return bToNDate(b, "2006-01-02", 4, 7)
}

// bToNMonthDayYear converts a byte slice to a date type (N days since Unix epoch, stored as int32)
// Follows a MM-DD-YYYY format; e.g., "01-02-2006"
func bToNMonthDayYear(b []byte) parsedRes {
	return bToNDate(b, "01-02-2006", 2, 5)
}

// bToNDayMonthYear converts a byte slice to a date type (N days since Unix epoch, stored as int32)
// Follows a DD-MM-YYYY format; e.g., "02-01-2006"
func bToNDayMonthYear(b []byte) parsedRes {
	return bToNDate(b, "02-01-2006", 2, 5)
}

// bToAMonthDayYearLong converts a byte slice to a date type (N days since Unix epoch, stored as int32)
// Follows a MM DD, YYYY format; e.g., "January 2, 2006"
func bToAMonthDayYearLong(b []byte) parsedRes {
	return bToADate(b, "January 2, 2006")
}

// bToAMonthDayYearShort converts a byte slice to a date type (N days since Unix epoch, stored as int32)
// Follows a mm DD, YYYY format; e.g., "Jan 2, 2006"
func bToAMonthDayYearShort(b []byte) parsedRes {
	return bToADate(b, "Jan 2, 2006")
}

// bToBool converts a byte slice to a boolean type
func bToBool(b []byte) parsedRes {
	var res boolRes = boolRes{val: false, isNull: true}

	const bytesInTrue = 4
	const bytesInFalse = 5

	switch len(b) {
	case 0:
		return res
	case 1:
		switch b[0] {
		case 'T', 't':
			res.val, res.isNull = true, false
		case 'F', 'f':
			res.val, res.isNull = false, false
		default:
			return res
		}
	case bytesInTrue:
		isTrue := (b[0] == 'T' || b[0] == 't') && b[1] == 'r' && b[2] == 'u' && b[3] == 'e'
		res.val, res.isNull = isTrue, isTrue != res.isNull
	case bytesInFalse:
		isFalse := (b[0] == 'F' || b[0] == 'f') && b[1] == 'a' && b[2] == 'l' && b[3] == 's' && b[4] == 'e'
		res.val, res.isNull = !isFalse, isFalse != res.isNull
	default:
		return res
	}
	return res
}

// bToSignedInteger converts a byte slice to signed integer type
func bToSignedInteger[T signedInteger](b []byte) numericRes[T] {
	var res numericRes[T] = numericRes[T]{val: 0, isNull: true}

	if len(b) == 0 {
		return res
	}

	var sign T = 1
	if (b[0] == dashChar) || (b[0] == plusChar) {
		if b[0] == dashChar {
			sign = -1
		}
		b = b[1:]
	}

	var val T = 0
	var base T = 10
	for i := 0; i < len(b); i++ {
		if !isNumericASCII(b[i]) {
			return res
		}
		val = val*base + T(b[i]-numericASCIILower)
	}

	res.val, res.isNull = val*sign, false

	return res
}

// bToFloatingPoint converts a byte slice to floating point type
func bToFloatingPoint[T floatingPoint](b []byte) numericRes[T] {
	var res numericRes[T] = numericRes[T]{val: 0, isNull: true}

	if len(b) == 0 {
		return res
	}

	var sign T = 1
	if (b[0] == dashChar) || (b[0] == plusChar) {
		if b[0] == dashChar {
			sign = -1
		}
		b = b[1:]
	}

	var whole T = 0
	var remainder T = 0
	var base T = 10
	var remBase T = 0.10
	var on int = 0

	for on < len(b) {
		if b[on] == decPntChar {
			on += 1
			break
		}
		if !isNumericASCII(b[on]) {
			return res
		}
		whole = whole*base + T(b[on]-numericASCIILower)
		on += 1
	}

	for on < len(b) {
		if !isNumericASCII(b[on]) {
			return res
		}
		remainder += remBase * T(b[on]-numericASCIILower)
		remBase *= 0.10
		on += 1
	}

	res.val, res.isNull = sign*(whole+remainder), false
	return res
}

type signedInteger interface {
	int32 | int64
}

type floatingPoint interface {
	float32 | float64
}

// bToNDate converts a byte slice to a date type
func bToNDate(b []byte, layout string, sepPos1, sepPos2 int) dateRes {
	var res dateRes = dateRes{val: 0, isNull: true}
	if len(b) == 0 {
		return res
	}

	b[sepPos1], b[sepPos2] = slashChar, slashChar

	d, err := time.Parse(layout, string(b))
	if err != nil {
		return res
	}

	const secsInOneDay int64 = 60 * 60 * 24
	res.val, res.isNull = int32(d.Unix()/secsInOneDay), false
	return res
}

// bToADate converts a byte slice to a date type
func bToADate(b []byte, layout string) dateRes {
	var res dateRes = dateRes{val: 0, isNull: true}
	if len(b) == 0 {
		return res
	}

	d, err := time.Parse(layout, string(b))
	if err != nil {
		return res
	}

	res.val, res.isNull = int32(d.Unix()), false
	return res
}
