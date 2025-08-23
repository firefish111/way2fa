package main

import (
	"fmt"
	"github.com/firefish111/way2fa/parse"
	"time"
)

func main() {
	csv, err := parse.GetDefaultCsv()
	if err != nil {
		panic(err)
	}

	accs, err := csv.GetAccs()
	if err != nil {
		panic(err)
	}

	curr := uint64(time.Now().Unix() / 30)

	for _, acc := range accs {
		code, err := acc.GenKey(curr)
		if err != nil {
			panic(err)
		}

		fmt.Printf("%s <@%s> : %d\n", acc.Name, acc.AcctId, code)
	}
}
