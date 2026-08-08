// Bubbletea bindings for Manager
package manager

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/firefish111/way2fa/detector"
	"github.com/firefish111/way2fa/internal/ui/common/styles"
	"github.com/firefish111/way2fa/parse"
)

func GetPossibilities() []parse.AccountList {
	ways := detector.GetPossibleWayFormats()
	pures := detector.GetPossiblePureFormats()

	ret := make([]parse.AccountList, 0, len(ways)+len(pures))
	for _, way := range ways {
		ret = append(ret, way)
	}
	for _, pure := range pures {
		ret = append(ret, pure)
	}

	return ret
}

type managerPhase int

const (
	selectStore managerPhase = iota
	setPassword
	waitingPassword
	done
)

type managerModel struct {
	possibilities []parse.AccountList
	selected      int
	destination   *parse.AccountList
	filename      *string
	phase         managerPhase
	helpModel     help.Model
	helpDB        map[string]key.Binding
}

func CreateCreatorModel(name *string, destination *parse.AccountList) *managerModel {
	ret := managerModel{
		selected:    0,
		destination: destination,
		filename:    name,
		phase:       selectStore,
		helpModel:   help.New(),
		helpDB:      defaultHelp(),
	}

	ret.helpModel.Styles.ShortDesc = styles.Faint
	ret.helpModel.Styles.FullDesc = styles.Faint

	ret.possibilities = GetPossibilities()

	return &ret
}
