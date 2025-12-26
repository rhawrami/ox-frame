package compute

import (
	"math"
	"sync"

	"github.com/rhawrami/ox-frame/ox/vector"
)

// AddLit returns the element-wise sum of a NumericVector of type T and literal value of type T
// AddLit will panic if both vectors are not of the same length
func AddLit[T vector.Numeric](x *vector.NumericVector[T], lit T) *vector.NumericVector[T] {
	return opLit(x, lit, addLitChunk)
}

// SubLit returns the element-wise difference of a NumericVector of type T and literal value of type T
//
// SubLit will panic if both vectors are not of the same length
func SubLit[T vector.Numeric](x *vector.NumericVector[T], lit T) *vector.NumericVector[T] {
	return opLit(x, lit, subLitChunk)
}

// MulLit returns the element-wise product of a NumericVector of type T and literal value of type T
//
// MulLit will panic if both vectors are not of the same length
func MulLit[T vector.Numeric](x *vector.NumericVector[T], lit T) *vector.NumericVector[T] {
	return opLit(x, lit, mulLitChunk)
}

// DivLit returns the element-wise quotient of a NumericVector of type T and literal value of type T
//
// DivLit will panic if both vectors are not of the same length; will also panic if divide by zero
func DivLit[T vector.Numeric](x *vector.NumericVector[T], lit T) *vector.NumericVector[T] {
	return opLit(x, lit, divLitChunk)
}

// PowLit returns the element-wise exponent expression of a NumericVector of type T and literal value of type T
func PowLit[T vector.Numeric](x *vector.NumericVector[T], lit T) *vector.NumericVector[T] {
	return opLit(x, lit, powLitChunk)
}

// opLit performs the vector operation between a vector and literal numeric, returning a resulting new vector.
func opLit[T vector.Numeric](x *vector.NumericVector[T], lit T, opFn func(out, x []T, lit T)) *vector.NumericVector[T] {
	dataBuff := make([]T, x.Len())
	validMap := x.Validity().DeepCopy()
	// break up chunks; make divisible by 8; final chunk will often not equal len of others
	chunkSize := x.Len() / (NumWorkers * 8) * 8

	xData := x.Data()

	var wg sync.WaitGroup
	wg.Add(NumWorkers)
	for i := 0; i < NumWorkers; i++ {
		// spawn workers
		go func(i int) {
			defer wg.Done()

			startData, endData := i*chunkSize, i*chunkSize+chunkSize
			// final chunk may not be div by 8
			if i == NumWorkers-1 {
				endData = x.Len()
			}
			// parallel vector operation
			opFn(
				dataBuff[startData:endData],
				xData[startData:endData],
				lit,
			)
		}(i)
	}
	wg.Wait()

	return vector.NumericVecFromComponents(
		x.Type(),
		dataBuff,
		validMap,
	)
}

// element wise vector scalar addition
func addLitChunk[T vector.Numeric](out, x []T, lit T) {
	for i := 0; i < len(out); i++ {
		out[i] = x[i] + lit
	}
}

// element wise vector scalar difference
func subLitChunk[T vector.Numeric](out, x []T, lit T) {
	for i := 0; i < len(out); i++ {
		out[i] = x[i] - lit
	}
}

// element wise vector scalar product
func mulLitChunk[T vector.Numeric](out, x []T, lit T) {
	for i := 0; i < len(out); i++ {
		out[i] = x[i] * lit
	}
}

// element wise vector scalar quotient
func divLitChunk[T vector.Numeric](out, x []T, lit T) {
	for i := 0; i < len(out); i++ {
		out[i] = x[i] / lit
	}
}

// element wise exponent expression
func powLitChunk[T vector.Numeric](out, x []T, lit T) {
	for i := 0; i < len(out); i++ {
		out[i] = T(math.Pow(float64(x[i]), float64(lit)))
	}
}
