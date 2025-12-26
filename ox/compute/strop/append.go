package strop

import (
	"sync"

	"github.com/rhawrami/ox-frame/ox/compute"
	"github.com/rhawrami/ox-frame/ox/vector"
)

// AddPrefix adds a prefix string (as byte slice input) to each string in a StringVector
func AddPrefix(x *vector.StringVector, s []byte) *vector.StringVector {
	return appendLit(x, s, addPrefix)
}

// AddSuffix adds a suffix string (as byte slice input) to each string in a StringVector
func AddSuffix(x *vector.StringVector, s []byte) *vector.StringVector {
	return appendLit(x, s, addSuffix)
}

func appendLit(x *vector.StringVector, lit []byte, opAppend func(cfg appendConfig)) *vector.StringVector {
	newLenB := len(x.Data()) + x.Len()*len(lit)

	newValidityMap := x.Validity().DeepCopy() // will stay same
	newOffsetsBuffer := make([]int64, len(x.Offsets()))
	newDataBuffer := make([]byte, newLenB)

	chunkSize := x.Len() / compute.NumWorkers

	var wg sync.WaitGroup
	wg.Add(compute.NumWorkers)
	for i := 0; i < compute.NumWorkers; i++ {
		go func(i int) {
			defer wg.Done()

			startIdx, stopIdx := i*chunkSize, i*chunkSize+chunkSize
			// handle final chunk
			if i == compute.NumWorkers-1 {
				stopIdx = x.Len()
			}

			cfg := appendConfig{
				dat0:  x.Data(),
				dat1:  newDataBuffer,
				off0:  x.Offsets(),
				off1:  newOffsetsBuffer,
				lit:   lit,
				start: startIdx,
				stop:  stopIdx,
			}

			opAppend(cfg)

		}(i)
	}
	wg.Wait()

	// handle final offset element
	newOffsetsBuffer[len(newOffsetsBuffer)-1] = int64(len(newDataBuffer))

	return vector.StringVecFromComponents(
		newDataBuffer,
		newOffsetsBuffer,
		newValidityMap,
	)
}

func addPrefix(cfg appendConfig) {
	for i := cfg.start; i < cfg.stop; i++ {
		newOffset := int(cfg.off0[i]) + (i-1)*len(cfg.lit)
		if i == 0 {
			newOffset = 0
		}
		newEnd := newOffset + int(cfg.off0[i+1]-cfg.off0[i]) + len(cfg.lit)

		cfg.off1[i] = int64(newOffset)

		copy(cfg.dat1[newOffset:newOffset+len(cfg.lit)], cfg.lit)
		copy(cfg.dat1[newOffset+len(cfg.lit):newEnd], cfg.dat0[cfg.off0[i]:cfg.off0[i+1]])
	}
}

func addSuffix(cfg appendConfig) {
	for i := cfg.start; i < cfg.stop; i++ {
		newOffset := int(cfg.off0[i]) + (i-1)*len(cfg.lit)
		if i == 0 {
			newOffset = 0
		}
		newEnd := newOffset + int(cfg.off0[i+1]-cfg.off0[i]) + len(cfg.lit)

		cfg.off1[i] = int64(newOffset)

		copy(cfg.dat1[newOffset:newEnd-len(cfg.lit)], cfg.dat0[cfg.off0[i]:cfg.off0[i+1]])
		copy(cfg.dat1[newEnd-len(cfg.lit):newEnd], cfg.lit)
	}
}

type appendConfig struct {
	dat0  []byte
	dat1  []byte
	off0  []int64
	off1  []int64
	lit   []byte
	start int
	stop  int
}
