package main

import (
	"fmt"

	"github.com/rhawrami/ox-frame/ox/compute"
	"github.com/rhawrami/ox-frame/ox/vector"
)

func main() {
	size := 100
	bools := make([]bool, size)
	strings := make([]string, size)

	for i := 0; i < size; i++ {
		bools[i] = true
		if i%13 == 0 {
			bools[i] = false
			strings[i] = ""
		} else {
			if i%6 == 0 {
				strings[i] = "Burckhardt"
			} else if i%10 == 0 {
				strings[i] = "Celine"
			} else if i%4 == 0 {
				strings[i] = "Bulgakov"
			} else if i%2 == 0 {
				strings[i] = "Dostoevsky"
			} else {
				strings[i] = "Soprano"
			}
		}
	}

	myStrVec := vector.StringVecFromStrings(strings, bools)
	for i := 0; i < myStrVec.Len(); i++ {
		fmt.Printf("%d\t%d\t[%s]\n", i, myStrVec.Offsets()[i], myStrVec.Data()[myStrVec.Offsets()[i]:myStrVec.Offsets()[i+1]])
	}

	s := []byte(" CODE999")

	myStrVecP := compute.AddPrefix(myStrVec, s)
	fmt.Println(myStrVecP.Len())
	myStrVecS := compute.AddSuffix(myStrVec, s)

	for i := 0; i < myStrVec.Len(); i++ {
		fmt.Printf("%d\n%s\n%s\n\n", myStrVecP.IsNullBinary(i), myStrVecP.StringValAt(i), myStrVecS.StringValAt(i))
	}
}
