package layouts

import (
	"gioui.org/layout"
)

type Flexing uint8

const (
	FirstFlexed Flexing = iota
	SecondFlexed
	FlexedRigid
	RigidFlexed
	DoubleRigid
)

type FlexWithTwoCildren struct {
	First, Second layout.Widget
	Flex          layout.Flex
	Weight        float32
	Flexing
}

func (flex FlexWithTwoCildren) Layout(gtx *layout.Context) {
	first := layout.Rigid(flex.First)
	second := layout.Rigid(flex.Second)

	switch flex.Flexing {
	case FirstFlexed:
		first = layout.Flexed(flex.Weight, flex.First)
		second = layout.Flexed(1-flex.Weight, flex.Second)
	case SecondFlexed:
		first = layout.Flexed(1-flex.Weight, flex.First)
		second = layout.Flexed(flex.Weight, flex.Second)
	case FlexedRigid:
		first = layout.Flexed(flex.Weight, flex.First)
	case RigidFlexed:
		second = layout.Flexed(flex.Weight, flex.Second)
	}
	flex.Flex.Layout(gtx, first, second)
}
func (flex FlexWithTwoCildren) Layedout(gtx *layout.Context) layout.Widget {
	return func() {
		flex.Layout(gtx)
	}
}
