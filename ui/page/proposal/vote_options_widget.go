package proposal

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
		increment: l.Theme.PlainIconButton(new(widget.Clickable), l.Icons.ContentAdd),
		decrement: l.Theme.PlainIconButton(new(widget.Clickable), l.Icons.ContentRemove),
		max:       l.Theme.Button(new(widget.Clickable), "MAX"),
	}
	i.max.Background = l.Theme.Color.Surface
	i.max.Color = l.Theme.Color.Gray2
	i.max.Font.Weight = text.Bold

	i.increment.Color, i.decrement.Color = l.Theme.Color.Text, l.Theme.Color.Text
	i.increment.Size, i.decrement.Size = values.TextSize18, values.TextSize18
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

	if i.max.Button.Clicked() {
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
