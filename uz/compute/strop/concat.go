package strop

import (
	"sync"

	"github.com/rhawrami/uz-frame/uz/compute"
	"github.com/rhawrami/uz-frame/uz/vector"
)

// Concat concatenates two StringVectors element-wise, with a byte slice separator
func Concat(x, y *vector.StringVector, sep []byte) *vector.StringVector {
	dataBuff := make([]byte, len(x.Data())+len(y.Data())+len(sep)*x.Len())
	offsetsBuff := make([]int64, len(x.Offsets()))
	validBuff := make([]byte, x.Validity().Len())

	chunkSize := x.Len() / compute.NumWorkers

	var wg sync.WaitGroup
	wg.Add(compute.NumWorkers)
	for i := 0; i < compute.NumWorkers; i++ {
		go func(i int) {
			defer wg.Done()

			startIdx, endIdx := i*chunkSize, i*chunkSize+chunkSize
			if i == compute.NumWorkers-1 {
				endIdx = x.Len()
			}

			cfg := &concatConfig{
				outData:    dataBuff,
				outValid:   validBuff,
				outOffsets: offsetsBuff,
				x:          x,
				y:          y,
				start:      startIdx,
				end:        endIdx,
				sep:        sep,
			}

			concatChunk(cfg)
		}(i)
	}
	wg.Wait()

	// handle final offset
	offsetsBuff[len(offsetsBuff)-1] = int64(len(dataBuff))

	newNullCount := vector.NullCountFromByteBuff(validBuff, x.Len())
	validity := vector.ValidityBitMap{
		TrueLen:   x.Len(),
		NullCount: newNullCount,
		Buffer:    validBuff,
	}

	return vector.StringVecFromComponents(
		dataBuff,
		offsetsBuff,
		validity,
	)
}

func concatChunk(cfg *concatConfig) {
	for i := cfg.start; i < cfg.end; i++ {
		newOffset := int(cfg.x.Offsets()[i]) + int(cfg.y.Offsets()[i]) + len(cfg.sep)*(i)
		xD, yD := cfg.x.ValAt(i), cfg.y.ValAt(i)

		cfg.outOffsets[i] = int64(newOffset)
		copy(cfg.outData[newOffset:newOffset+len(xD)], xD)
		copy(cfg.outData[newOffset+len(xD):newOffset+len(xD)+len(cfg.sep)], cfg.sep)
		copy(cfg.outData[newOffset+len(xD)+len(cfg.sep):newOffset+len(xD)+len(cfg.sep)+len(yD)], yD)

		if i%8 == 0 {
			cfg.outValid[i/8] = cfg.x.Validity().Buffer[i/8] & cfg.y.Validity().Buffer[i/8]
		}
	}
}

type concatConfig struct {
	outData    []byte
	outValid   []byte
	outOffsets []int64
	x          *vector.StringVector
	y          *vector.StringVector
	start      int
	end        int
	sep        []byte
}
