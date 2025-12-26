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

func ToUpperASCII(x *vector.StringVector) *vector.StringVector {
	return changeCaseASCII(x, toUpperASCII)
}

func ToLowerASCII(x *vector.StringVector) *vector.StringVector {
	return changeCaseASCII(x, toLowerASCII)
}

func SwapCaseASCII(x *vector.StringVector) *vector.StringVector {
	return changeCaseASCII(x, swapCaseASCII)
}

func changeCaseASCII(x *vector.StringVector, opFn func(cfg caseChangeConfig)) *vector.StringVector {
	newVec := x.DeepCopy()

	chunkSize := len(x.Data()) / compute.NumWorkers

	var wg sync.WaitGroup
	wg.Add(compute.NumWorkers)
	for i := 0; i < compute.NumWorkers; i++ {
		go func(i int) {
			defer wg.Done()

			startData, endData := i*chunkSize, i*chunkSize+chunkSize
			// handle final chunk
			if i == compute.NumWorkers-1 {
				endData = len(x.Data())
			}

			cfg := caseChangeConfig{
				outData: newVec.Data(),
				xData: x.Data(),
				start: ,
			}
			opFn(newVec.Data()[startData:endData])
		}(i)
	}
	wg.Wait()

	return newVec
}

type caseChangeConfig struct {
	outData []byte
	xData   []byte
	start   int
	stop    int
}

func toUpperASCII(cfg *caseChangeConfig) {
	in := cfg.xData
	out := cfg.outData
	for i := 0; i < len(in); i++ {
		if isLowerASCII(in[i]) {
			out[i] = in[i] - diffLowerUpper
		}
	}
}

func toLowerASCII(cfg *caseChangeConfig) {
	in := cfg.xData
	out := cfg.outData
	for i := 0; i < len(in); i++ {
		if isLowerASCII(in[i]) {
			out[i] = in[i] + diffLowerUpper
		}
	}
}

func swapCaseASCII(cfg *caseChangeConfig) {
	in := cfg.xData
	out := cfg.outData
	for i := 0; i < len(in); i++ {
		ogVal := in[i]
		if isLowerASCII(ogVal) {
			out[i] = in[i] - diffLowerUpper
			continue
		}
		if isUpperASCII(ogVal) {
			out[i] = in[i] + diffLowerUpper
		}
	}
}

func isLowerASCII(x byte) bool {
	return x >= asciiLowerCaseMin && x <= asciiLowerCaseMax
}

func isUpperASCII(x byte) bool {
	return x >= asciiUpperCaseMin && x <= asciiUpperCaseMax
}
