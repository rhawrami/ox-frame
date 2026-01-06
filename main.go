package main

import (
	"fmt"

	"github.com/rhawrami/uz-frame/uz/io/csv"
)

func main() {
	b, err := csv.GetHeader("test.csv", ',', '\n')

	if err != nil {
		fmt.Println(err)
		return
	}

	for _, v := range b {
		fmt.Println(string(v))
	}
}
