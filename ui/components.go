package ui

import (
	"fmt"

	"gioui.org/layout"
	//"gioui.org/unit"
	"github.com/atotto/clipboard"
	//"github.com/raedahgroup/godcr/ui/decredmaterial"
)

const (
	headerHeight = .15
	navSize      = .1
)

var (
	// layout.Flex: Vertical
	vertFlex = layout.Flex{Axis: layout.Vertical}
	// layout.Flex: Horizontal
	horFlex = layout.Flex{}
	// layout.Rigid
	rigid = layout.Rigid
)

func toMax(gtx *layout.Context) {
	gtx.Constraints.Width.Min = gtx.Constraints.Width.Max
	gtx.Constraints.Height.Min = gtx.Constraints.Height.Max
}

func GetClipboardContent() string {
	str, err := clipboard.ReadAll()
	if err != nil {
		log.Warn(fmt.Sprintf("error getting clipboard data: %s", err.Error()))
		return ""
	}

	return str
}
