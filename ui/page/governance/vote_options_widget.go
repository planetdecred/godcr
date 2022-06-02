package governance

import (
	"fmt"
	"image/color"
	"strconv"

	"gioui.org/text"
	"gioui.org/widget"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/values"
)

type inputVoteOptionsWidgets struct {
	label     string
	activeBg  color.NRGBA
	dotColor  color.NRGBA
	input     decredmaterial.Editor
	increment decredmaterial.IconButton
	decrement decredmaterial.IconButton
	max       decredmaterial.Button
}

func newInputVoteOptions(l *load.Load, label string) *inputVoteOptionsWidgets {
	i := &inputVoteOptionsWidgets{
		label:     label,
		activeBg:  l.Theme.Color.Green50,
		dotColor:  l.Theme.Color.Green500,
		input:     l.Theme.Editor(new(widget.Editor), ""),
		increment: l.Theme.IconButton(l.Theme.Icons.ContentAdd),
		decrement: l.Theme.IconButton(l.Theme.Icons.ContentRemove),
		max:       l.Theme.Button(values.String(values.StrMax)),
	}
	i.max.Background = l.Theme.Color.Surface
	i.max.Color = l.Theme.Color.GrayText1
	i.max.Font.Weight = text.SemiBold

	i.increment.ChangeColorStyle(&values.ColorStyle{Foreground: l.Theme.Color.Text})
	i.decrement.ChangeColorStyle(&values.ColorStyle{Foreground: l.Theme.Color.Text})

	i.increment.Size, i.decrement.Size = values.MarginPadding18, values.MarginPadding18
	i.input.Bordered = false
	i.input.Editor.SetText("0")
	i.input.Editor.Alignment = text.Middle
	return i
}

func (i *inputVoteOptionsWidgets) voteCount() int {
	value, err := strconv.Atoi(i.input.Editor.Text())
	if err != nil {
		fmt.Println(err)
		return 0
	}
	return value
}

func (i *inputVoteOptionsWidgets) reset() {
	i.input.Editor.SetText("0")
}

func (vm *voteModal) handleVoteCountButtons(i *inputVoteOptionsWidgets) {
	if i.increment.Button.Clicked() {
		value, err := strconv.Atoi(i.input.Editor.Text())
		if err != nil {
			return
		}
		if vm.remainingVotes() <= 0 {
			return
		}
		value++
		i.input.Editor.SetText(fmt.Sprintf("%d", value))
	}

	if i.decrement.Button.Clicked() {
		value, err := strconv.Atoi(i.input.Editor.Text())
		if err != nil {
			return
		}
		value--
		if value < 0 {
			return
		}
		i.input.Editor.SetText(fmt.Sprintf("%d", value))
	}

	if i.max.Clicked() {
		max := vm.remainingVotes() + i.voteCount()
		i.input.Editor.SetText(fmt.Sprint(max))
	}

	for _, e := range i.input.Editor.Events() {
		switch e.(type) {
		case widget.ChangeEvent:
			count := i.voteCount()
			if count < 0 {
				i.input.Editor.SetText("0")
			}
		}
	}
}
