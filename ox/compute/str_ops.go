package compute

import (
	"sync"

	"github.com/rhawrami/ox-frame/ox/vector"
)

func AppendLit(x *vector.StringVector, lit string, opAppend func(outDat, xDat []byte, outOff, xOff []int64, lit []byte)) *vector.StringVector {
	litB := []byte(lit)
	newLenB := len(x.Data()) + (len(x.Data()) * len(litB))

	newValidityMapBuffer := x.Validity().DeepCopyBuff() // will stay same
	newOffsetsBuffer := make([]int64, len(x.Offsets()))
	newDataBuffer := make([]byte, newLenB)

	chunkSize := x.Len()/NumWorkers + 1
	if x.Len()%NumWorkers == 0 {
		chunkSize = chunkSize - 1
	}

	var wg sync.WaitGroup
	wg.Add(NumWorkers)
	for i := 0; i < NumWorkers; i++ {
		go func(i int) {
			defer wg.Done()

			startElement, endElement := i*chunkSize, i*chunkSize+chunkSize

			opAppend(
				newDataBuffer,
				x.Data(),
				newOffsetsBuffer[startElement:endElement],
				x.Offsets()[startElement:endElement],
				litB,
			)

		}(i)
	}
	wg.Wait()

	return &vector.StringVector{}
}

func appendRight(outD, xD []byte, outO, xO []int64, lit []byte) {
	for i := 0; i < len(outO); i++ {
		newOffSet := 0
	}
}
