package manager

import "fmt"

// list all available stores
func List() {
	fmt.Println("Available default stores, listed in order of priority.\n\n" + disclaimer)
}
