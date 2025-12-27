package strop

import (
	"sync"

	"github.com/rhawrami/ox-frame/ox/compute"
	"github.com/rhawrami/ox-frame/ox/vector"
)

const (
	diffLowerUpper    byte = 32
	asciiUpperCaseMin byte = 65
	asciiUpperCaseMax byte = 90
	asciiLowerCaseMin byte = 97
	asciiLowerCaseMax byte = 122
)

// ToTitleASCII converts elements in a StringVector to title-case; assumes bytes are ASCII
func ToTitleASCII(x *vector.StringVector) *vector.StringVector {
	return changeCaseASCII(x, toTitleASCII)
}

// ToUpperASCII converts elements in a StringVector to upper-case; assumes bytes are ASCII
func ToUpperASCII(x *vector.StringVector) *vector.StringVector {
	return changeCaseASCII(x, toUpperASCII)
}

// ToLowerASCII converts elements in a StringVector to lower-case; assumes bytes are ASCII
func ToLowerASCII(x *vector.StringVector) *vector.StringVector {
	return changeCaseASCII(x, toLowerASCII)
}

// SwapCaseASCII changes the case of each alphabetical byte in a StringVector; assumes bytes are ASCII
func SwapCaseASCII(x *vector.StringVector) *vector.StringVector {
	return changeCaseASCII(x, swapCaseASCII)
}

func changeCaseASCII(x *vector.StringVector, opFn func(x *vector.StringVector, start, stop int)) *vector.StringVector {
	newVec := x.DeepCopy()

	chunkSize := newVec.Len() / compute.NumWorkers

	var wg sync.WaitGroup
	wg.Add(compute.NumWorkers)
	for i := 0; i < compute.NumWorkers; i++ {
		go func(i int) {
			defer wg.Done()

			startIdx, endIdx := i*chunkSize, i*chunkSize+chunkSize
			// handle final chunk
			if i == compute.NumWorkers-1 {
				endIdx = newVec.Len()
			}
			opFn(newVec, startIdx, endIdx)
		}(i)
	}
	wg.Wait()

	return newVec
}

func toTitleASCII(x *vector.StringVector, start, stop int) {
	for i := start; i < stop; i++ {
		toTitleASCIIOneWord(x.Data()[x.Offsets()[i]:x.Offsets()[i+1]])
	}
}

func toTitleASCIIOneWord(x []byte) {
	// first byte to upper
	if len(x) > 0 && isLowerASCII(x[0]) {
		x[0] = x[0] - diffLowerUpper
	}
	if len(x) <= 1 {
		return
	}
	for i := 1; i < len(x); i++ {
		if isUpperASCII(x[i]) {
			x[i] = x[i] + diffLowerUpper
		}
	}
}

func toUpperASCII(x *vector.StringVector, start, stop int) {
	for i := x.Offsets()[start]; i < x.Offsets()[stop]; i++ {
		if isLowerASCII(x.Data()[i]) {
			x.Data()[i] = x.Data()[i] - diffLowerUpper
		}
	}
}

func toLowerASCII(x *vector.StringVector, start, stop int) {
	for i := x.Offsets()[start]; i < x.Offsets()[stop]; i++ {
		if isUpperASCII(x.Data()[i]) {
			x.Data()[i] = x.Data()[i] + diffLowerUpper
		}
	}
}

func swapCaseASCII(x *vector.StringVector, start, stop int) {
	for i := x.Offsets()[start]; i < x.Offsets()[stop]; i++ {
		ogVal := x.Data()[i]
		if isLowerASCII(ogVal) {
			x.Data()[i] = x.Data()[i] - diffLowerUpper
			continue
		}
		if isUpperASCII(ogVal) {
			x.Data()[i] = x.Data()[i] + diffLowerUpper
		}
	}
}

func isLowerASCII(x byte) bool {
	return x >= asciiLowerCaseMin && x <= asciiLowerCaseMax
}

func isUpperASCII(x byte) bool {
	return x >= asciiUpperCaseMin && x <= asciiUpperCaseMax
}
