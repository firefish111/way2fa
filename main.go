package main

import (
	"flag"
	"fmt"
	"github.com/firefish111/way2fa/parse/csv"
	"time"
)

const (
	VersionMajor = 0
	VersionMinor = 2
	VersionPatch = 0
)

func main() {
	ver := flag.Bool("version", false, "Prints version")

	flag.Parse()

	// end of args

	if *ver {
		fmt.Printf("way2fa v%d.%d.%d\n", VersionMajor, VersionMinor, VersionPatch)
		return
	}

	var name *string
	if a := flag.Args(); len(a) > 0 {
		name = &a[0]
	}

	store, err := csv.GetFile(name)
	if err != nil {
		panic(err)
	}

	accs, err := store.GetAccs()
	if err != nil {
		panic(err)
	}

	curr_time := uint64(time.Now().Unix())

	srct, srcs := store.GetSource()
	fmt.Printf("%d: %s\n", srct, srcs)

	for _, acc := range accs {
		code, err := acc.GenKey(curr_time / uint64(acc.GetInterval()))
		if err != nil {
			panic(err)
		}

		fmt.Printf("%s <@%s> : %d\n", acc.Name, acc.AcctId, code)
	}
}
