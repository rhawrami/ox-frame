package csv

const (
	dashChar   byte = '-'
	slashChar  byte = '/'
	spaceChar  byte = ' '
	decPntChar byte = '.'
	percChar   byte = '%'
	ucEChar    byte = 'E'
	lcEChar    byte = 'e'

	numericASCIILower byte = '0'
	numericASCIIUpper byte = '9'

	alphaASCIIUCMin    byte = 'A'
	alphaASCIIUCMax    byte = 'Z'
	alphaASCIILCMin    byte = 'a'
	alphaASCIILCMax    byte = 'z'
	alphaASCIIUCLCDIff byte = 32
)

type inferredType int

const (
	// null value (e.g., empty string)
	null inferredType = iota

	// numeric
	notNum
	intNum   // 10
	floatNum // 10.0

	// date
	notDate
	nYearMonthDay // 2006-01-02
	nMonthDayYear // 01-02-2006
	nDayMonthYear // 02-01-2006
	aMonthDayYear // [Jj]an(uary).? 2(nd)?,? 2006

	// boolean
	notBoolean
	boolean // true

	// string (default)
	text
)

func IsNumeric(b []byte) inferredType {
	if len(b) == 0 {
		return null
	}
	if b[0] == dashChar {
		b = b[1:]
	}

	foundOther := false
	foundNum := false
	foundDecPnt := false

	for i, v := range b {
		if v >= numericASCIILower && v <= numericASCIIUpper {
			foundNum = true
			continue
		}
		switch v {
		case decPntChar:
			foundDecPnt = true
		case percChar:
			// percentage char should be at the end
			foundOther = i != (len(b) - 1)
		case ucEChar, lcEChar:
			// e|E cannot be at start or end
			foundOther = (i == (len(b) - 1)) || (i == 0)
		default:
			foundOther = true
		}
	}
	if foundNum && !foundOther {
		if foundDecPnt {
			return floatNum
		}
		return intNum
	}
	return notNum
}

func IsDate(b []byte) inferredType {
	if len(b) == 0 {
		return null
	}
	// minimum length of the byte slice should be 10
	// shortest possible date strings:
	// - e.g., `01-02-2006` (len 10)
	// - e.g., `jan 1 2006` (len 10)
	// maximum length should be 20
	// - e.g., `September 30th, 2026`
	minLen := 10
	maxLen := 20
	if len(b) < minLen || len(b) > maxLen {
		return notDate
	}

	// if like \d{2,4}[-/]\d{2}[-/]\d{2,4}, b[2] or b[4] must be [-/]
	// must also be len 10
	if (len(b) == minLen) && (b[2] == slashChar || b[4] == slashChar || b[2] == dashChar || b[4] == dashChar) {
		return isNDate(b)
	}

	// if like [Jj]an(uary).? 2(nd)?,? 2006, b[0] must be ASCII
	if isAlphaASCII(b[0]) {
		return isADate(b)
	}

	return notDate
}

func isNDate(b []byte) inferredType {
	// if sep is on b[4], must be like `2006-01-02`
	if b[4] == slashChar || b[4] == dashChar {
		// ensure that b[7] is also sep
		if b[7] == slashChar || b[7] == dashChar {
			return nYearMonthDay
		}
		return notDate
	}
	// when sep is on b[2], need to figure out if
	// `01-02-2006` (MonthDayYear) OR
	// `02-01-2006` (DayMonthYear)
	// check b[0] and b[1]
	// IF MonthDayYear
	// - b[0] must ~ [01]
	if b[0] != '0' && b[0] != '1' {
		return nDayMonthYear
	}

	if b[0] == '1' && b[1] != '0' && b[1] != '1' && b[1] != '2' {
		return nDayMonthYear
	}

	return nMonthDayYear
}

func isADate(b []byte) inferredType {
	// ensure that last 4 bytes are numeric
	for i := 0; i < 4; i++ {
		if !isNumericASCII(b[len(b)-i-1]) {
			return notDate
		}
	}
	return aMonthDayYear
}

func IsBool(b []byte) inferredType {
	var bT inferredType

	switch len(b) {
	case 0:
		bT = null
	case 1, 4, 5:
		// first char
		if b[0] == 'f' || b[0] == 't' || b[0] == 'F' || b[0] == 'T' {
			if len(b) == 1 {
				bT = boolean
			}
			// `true` and `false` share last char
			if b[len(b)-1] == 'e' || b[len(b)-1] == 'E' {
				switch len(b) {
				case 4:
					// check `true`
					if b[1] == 'r' || b[1] == 'R' {
						if b[2] == 'u' || b[2] == 'U' {
							bT = boolean
						}
					}
				case 5:
					// check `false`
					if b[1] == 'a' || b[1] == 'A' {
						if b[2] == 'l' || b[2] == 'L' {
							if b[3] == 's' || b[3] == 'S' {
								bT = boolean
							}
						}
					}
				default:
				}
			}
		}
	default:
		bT = notBoolean
	}
	return bT
}

func isAlphaASCII(b byte) bool {
	return (b >= alphaASCIIUCMin && b <= alphaASCIIUCMax) || (b >= alphaASCIILCMin && b <= alphaASCIILCMax)
}

func isNumericASCII(b byte) bool {
	return b >= numericASCIILower && b <= numericASCIIUpper
}
