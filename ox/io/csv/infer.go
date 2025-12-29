package csv

const (
	dashChar   byte = '-'
	plusChar   byte = '+'
	slashChar  byte = '/'
	spaceChar  byte = ' '
	decPntChar byte = '.'
	percChar   byte = '%'
	ucEChar    byte = 'E'
	lcEChar    byte = 'e'

	numericASCIILower byte = '0'
	numericASCIIUpper byte = '9'

	alphaASCIIUCMin byte = 'A'
	alphaASCIIUCMax byte = 'Z'
	alphaASCIILCMin byte = 'a'
	alphaASCIILCMax byte = 'z'
)

func newColInferrer(cName string) *colInferrer {
	tally := newInferenceTally()
	return &colInferrer{
		cName: cName,
		tally: tally,
	}
}

// colInferrer manages inference from a column in a CSV
type colInferrer struct {
	cName     string
	tally     inferenceTally
	valLenMin int
	valLenMax int
	valLenSum int
	nSample   int
}

func (c *colInferrer) predictParser() func([]byte) parsedRes {
	switch c.predictType() {
	case intNum:
		const i32MaxLen = 10
		if c.valLenMax < (i32MaxLen) {
			return bToInt32
		}
		return bToInt64

	case floatNum:
		// come back for float32 vs float64 impl
		return bToFloat64

	case nYearMonthDay:
		return bToNYearMonthDay
	case nMonthDayYear:
		return bToNMonthDayYear
	case nDayMonthYear:
		return bToNDayMonthYear
	case aMonthDayYearLong:
		return bToAMonthDayYearLong
	case aMonthDayYearShort:
		return bToAMonthDayYearShort

	case boolean:
		return bToBool

	case strDefault:
		return bToStr
	}

	return bToStr
}

// predictType makes a final prediction for the type of a column
func (c *colInferrer) predictType() inferredType {
	var (
		firstPred, secondPred           inferredType = strDefault, strDefault
		firstPredShare, secondPredShare float32      = 0, 0
	)

	// get shares
	c.tally.normalize(float32(c.nSample))

	for k, kShare := range c.tally {
		if kShare > secondPredShare {
			if kShare > firstPredShare {
				prevFirstPred := firstPred
				prevFirstPredShare := firstPredShare
				firstPred = k
				firstPredShare = kShare
				if prevFirstPredShare > secondPredShare {
					secondPred = prevFirstPred
					secondPredShare = prevFirstPredShare
				}
				continue
			}
			secondPred = k
			secondPredShare = kShare
		}
	}

	const SeventyFivePerc float32 = 0.75
	if firstPredShare >= SeventyFivePerc && firstPred != null {
		return firstPred
	}

	const FiftyPerc, TwentyFivePerc float32 = 0.50, 0.25
	if firstPredShare >= FiftyPerc && firstPred != null {
		if secondPred != null && secondPredShare >= TwentyFivePerc {
			// deal with common cases

			// numeric mixture -> use float
			if (firstPred == floatNum && secondPred == intNum) || (firstPred == intNum && secondPred == floatNum) {
				return floatNum
			}
		}
		return firstPred
	}

	// default to string
	return strDefault
}

func (c *colInferrer) updateStatistics(b []byte) {
	// incr sample size
	c.nSample += 1
	// update inferenceTally
	t := inferType(b)
	c.tally.updateTally(t)

	// value length statistics
	c.updateValLenStatistics(b)
}

func (c *colInferrer) updateValLenStatistics(b []byte) {
	if c.nSample == 1 {
		c.valLenSum = len(b)
		c.valLenMin = len(b)
		c.valLenMax = len(b)
		return
	}
	if len(b) < c.valLenMin {
		c.valLenMin = len(b)
	}
	if len(b) > c.valLenMax {
		c.valLenMax = len(b)
	}
	c.valLenSum += len(b)
}

// inferType infers the likely type of an input
func inferType(b []byte) inferredType {
	// in order, check:
	// Null (e.g., len 0)
	// Numeric
	// Date
	// Boolean
	// String as default

	if len(b) == 0 {
		return null
	}

	switch isNumeric(b) {
	case intNum:
		return intNum
	case floatNum:
		return floatNum
	default:
		//
	}

	switch isDate(b) {
	case nYearMonthDay:
		return nYearMonthDay
	case nMonthDayYear:
		return nMonthDayYear
	case nDayMonthYear:
		return nDayMonthYear
	case aMonthDayYearLong:
		return aMonthDayYearLong
	case aMonthDayYearShort:
		return aMonthDayYearShort
	default:
		//
	}

	switch isBool(b) {
	case boolean:
		return boolean
	default:
		//
	}

	return strDefault
}

// inferenceTally keeps a tally of inferred type for a column
type inferenceTally map[inferredType]float32

func (t inferenceTally) normalize(n float32) {
	for k, v := range t {
		t[k] = v / n
	}
}

func (t inferenceTally) updateTally(i inferredType) {
	t[i] += 1
}

func newInferenceTally() inferenceTally {
	return inferenceTally{
		null:               0,
		intNum:             0,
		floatNum:           0,
		nYearMonthDay:      0,
		nMonthDayYear:      0,
		nDayMonthYear:      0,
		aMonthDayYearLong:  0,
		aMonthDayYearShort: 0,
		boolean:            0,
		strDefault:         0,
	}
}

type inferredType int

const (
	// null value (e.g., empty slice)
	null inferredType = iota

	// numeric
	notNum
	intNum   // 10
	floatNum // 10.0

	// date
	notDate
	nYearMonthDay      // 2006-01-02
	nMonthDayYear      // 01-02-2006
	nDayMonthYear      // 02-01-2006
	aMonthDayYearLong  // January 2, 2006
	aMonthDayYearShort // Jan 2, 2006

	// boolean
	notBoolean
	boolean // true

	// string (default)
	strDefault
)

// isNumeric determines if an input is likely a numeric type
func isNumeric(b []byte) inferredType {
	if b[0] == dashChar {
		b = b[1:]
	}

	foundOther := false
	foundNum := false
	foundDecPnt := false

	for i, v := range b {
		if isNumericASCII(v) {
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

// isDate determines if an input is likely a date type
func isDate(b []byte) inferredType {
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

	// check byte 3 for space or character
	// in case of abbreviated month, byte 3 will be ' ', EXCEPT FOR `May`
	if b[3] == spaceChar && b[3] != 'y' {
		return aMonthDayYearShort
	}

	return aMonthDayYearLong
}

// isBool determines if an input is likely a boolean type
func isBool(b []byte) inferredType {
	var bT inferredType

	switch len(b) {
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
