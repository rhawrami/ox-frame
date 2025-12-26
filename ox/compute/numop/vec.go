package numop

import (
	"sync"

	"github.com/rhawrami/ox-frame/ox/compute"
	"github.com/rhawrami/ox-frame/ox/vector"
)

// AddVec returns the element-wise sum of two NumericVectors of type T
//
// AddVec will panic if both vectors are not of the same length
func AddVec[T vector.Numeric](x, y *vector.NumericVector[T]) *vector.NumericVector[T] {
	return opVec(x, y, addVecChunk8Incr)
}

// SubVec returns the element-wise difference of two NumericVectors of type T
//
// SubVec will panic if both vectors are not of the same length
func SubVec[T vector.Numeric](x, y *vector.NumericVector[T]) *vector.NumericVector[T] {
	return opVec(x, y, subVecChunk8Incr)
}

// MulVec returns the element-wise product of two NumericVectors of type T
//
// MulVec will panic if both vectors are not of the same length
func MulVec[T vector.Numeric](x, y *vector.NumericVector[T]) *vector.NumericVector[T] {
	return opVec(x, y, mulVecChunk8Incr)
}

// opVec performs the binary vector operation on two vectors, returning a resulting new vector.
func opVec[T vector.Numeric](x, y *vector.NumericVector[T], opFn func(out, x, y []T, outB, xB, yB []byte)) *vector.NumericVector[T] {
	dataBuff := make([]T, x.Len())
	validBuff := make([]byte, x.Validity().Len())
	// break up chunks; make divisible by 8; final chunk will often not equal len of others
	chunkSize := x.Len() / (compute.NumWorkers * 8) * 8

	xData, yData := x.Data(), y.Data()
	xValid, yValid := x.Validity().Buffer, y.Validity().Buffer

	var wg sync.WaitGroup
	wg.Add(compute.NumWorkers)
	for i := 0; i < compute.NumWorkers; i++ {
		// spawn workers
		go func(i int) {
			defer wg.Done()

			startData, endData := i*chunkSize, i*chunkSize+chunkSize
			startValidity, endValidity := i*chunkSize/8, (i*chunkSize+chunkSize)/8
			// final chunk may not be div by 8
			if i == compute.NumWorkers-1 {
				endData = x.Len()
				endValidity = x.Validity().Len()
			}
			// parallel vector operation
			opFn(
				dataBuff[startData:endData],
				xData[startData:endData],
				yData[startData:endData],
				validBuff[startValidity:endValidity],
				xValid[startValidity:endValidity],
				yValid[startValidity:endValidity],
			)
		}(i)
	}
	wg.Wait()

	// new validMap
	nullCount := vector.NullCountFromByteBuff(validBuff, x.Len())
	validMap := vector.ValidityBitMap{
		TrueLen:   x.Len(),
		NullCount: nullCount,
		Buffer:    validBuff,
	}

	return vector.NumericVecFromComponents(
		x.Type(),
		dataBuff,
		validMap,
	)
}

// element-wise vector sum
func addVecChunk8Incr[T vector.Numeric](out, x, y []T, outB, xB, yB []byte) {
	for i := 0; i < len(out); i += 8 {
		// unroll 8 elems
		out[i] = x[i] + y[i]
		out[i+1] = x[i+1] + y[i+1]
		out[i+2] = x[i+2] + y[i+2]
		out[i+3] = x[i+3] + y[i+3]
		out[i+4] = x[i+4] + y[i+4]
		out[i+5] = x[i+5] + y[i+5]
		out[i+6] = x[i+6] + y[i+6]
		out[i+7] = x[i+7] + y[i+7]
		// bitwise AND to get new nulls
		outB[i/8] = xB[i/8] & yB[i/8]
	}
	// if working on last chunk, evaluate remainders
	if len(out)%8 != 0 {
		// final data points
		for i := len(out) / 8 * 8; i < len(out); i++ {
			out[i] = x[i] + y[i]
		}
		// final validity byte
		outB[len(outB)-1] = xB[len(xB)-1] & yB[len(yB)-1]
	}
}

// element-wise vector difference
func subVecChunk8Incr[T vector.Numeric](out, x, y []T, outB, xB, yB []byte) {
	for i := 0; i < len(out); i += 8 {
		// unroll 8 elems
		out[i] = x[i] - y[i]
		out[i+1] = x[i+1] - y[i+1]
		out[i+2] = x[i+2] - y[i+2]
		out[i+3] = x[i+3] - y[i+3]
		out[i+4] = x[i+4] - y[i+4]
		out[i+5] = x[i+5] - y[i+5]
		out[i+6] = x[i+6] - y[i+6]
		out[i+7] = x[i+7] - y[i+7]
		// bitwise AND to get new nulls
		outB[i/8] = xB[i/8] & yB[i/8]
	}
	// if working on last chunk, evaluate remainders
	if len(out)%8 != 0 {
		// final data points
		for i := len(out) / 8 * 8; i < len(out); i++ {
			out[i] = x[i] - y[i]
		}
		// final validity byte
		outB[len(outB)-1] = xB[len(xB)-1] & yB[len(yB)-1]
	}
}

// element-wise vector product
func mulVecChunk8Incr[T vector.Numeric](out, x, y []T, outB, xB, yB []byte) {
	for i := 0; i < len(out); i += 8 {
		// unroll 8 elems
		out[i] = x[i] * y[i]
		out[i+1] = x[i+1] * y[i+1]
		out[i+2] = x[i+2] * y[i+2]
		out[i+3] = x[i+3] * y[i+3]
		out[i+4] = x[i+4] * y[i+4]
		out[i+5] = x[i+5] * y[i+5]
		out[i+6] = x[i+6] * y[i+6]
		out[i+7] = x[i+7] * y[i+7]
		// bitwise AND to get new nulls
		outB[i/8] = xB[i/8] & yB[i/8]
	}
	// if working on last chunk, evaluate remainders
	if len(out)%8 != 0 {
		// final data points
		for i := len(out) / 8 * 8; i < len(out); i++ {
			out[i] = x[i] * y[i]
		}
		// final validity byte
		outB[len(outB)-1] = xB[len(xB)-1] & yB[len(yB)-1]
	}
}
