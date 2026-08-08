package manager

import (
	"fmt"

	"github.com/firefish111/way2fa/internal/ui/common/styles"
	managerUi "github.com/firefish111/way2fa/internal/ui/manager"
)

// list all available stores
func List() {
	fmt.Println(styles.SidePad.Render("\nAvailable default stores, listed in order of priority.\nFor .way types, its corresponding ID number is also shown."))

	tab := managerUi.MakeListTable(nil)
	fmt.Println(tab)
}
