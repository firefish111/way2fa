// Bubbletea bindings for Manager
package manager

import (
	"github.com/firefish111/way2fa/detector"
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
	selected      parse.AccountList
	destination   *parse.AccountList
	filename      *string
	phase         managerPhase
}

func CreateCreatorModel(name *string, destination *parse.AccountList) *managerModel {
	ret := managerModel{
		selected:    nil,
		destination: destination,
		filename:    name,
		phase:       selectStore,
	}

	ret.possibilities = GetPossibilities()

	return &ret
}
