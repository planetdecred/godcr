package modal

import (
	"math/rand"
	"time"

	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/widget"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

func editorsNotEmpty(editors ...*widget.Editor) bool {
	for _, e := range editors {
		if e.Text() == "" {
			return false
		}
	}
	return true
}

func generateRandomNumber() int {
	return rand.New(rand.NewSource(time.Now().UnixNano())).Int()
}

func computePasswordStrength(pb *decredmaterial.ProgressBarStyle, th *decredmaterial.Theme, editors ...*widget.Editor) {
	password := editors[0]
	strength := dcrlibwallet.ShannonEntropy(password.Text()) / 4.0
	pb.Progress = float32(strength)

	//set progress bar color
	switch {
	case pb.Progress <= 0.30:
		pb.Color = th.Color.Danger
	case pb.Progress > 0.30 && pb.Progress <= 0.60:
		pb.Color = th.Color.Yellow
	case pb.Progress > 0.50:
		pb.Color = th.Color.Success
	}

}

func handleTabEvent(event chan *key.Event) bool {
	var isTabPressed bool
	select {
	case event := <-event:
		if event.Name == key.NameTab && event.State == key.Press {
			isTabPressed = true
		}
	default:
	}
	return isTabPressed
}

func SwitchEditors(keyEvent chan *key.Event, editors ...*widget.Editor) {
	for i := 0; i < len(editors); i++ {
		if editors[i].Focused() {
			if handleTabEvent(keyEvent) {
				if i == len(editors)-1 {
					editors[0].Focus()
				} else {
					editors[i+1].Focus()
				}
			}
		}
	}
}
