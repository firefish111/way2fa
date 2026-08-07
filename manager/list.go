package manager

import (
	"fmt"
	"reflect"

	"github.com/firefish111/way2fa/internal/ui/manager"
)

// list all available stores
func List() {
	fmt.Println("Available default stores, listed in order of priority.\n\n" + disclaimer)

	for i, store := range manager.GetPossibilities() {
		fmt.Printf("\t[%d]: %s\n", i+1, reflect.TypeOf(store).Elem().Name())
	}
}
