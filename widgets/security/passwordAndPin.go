package security

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr-gio/helper"
	"github.com/raedahgroup/godcr-gio/widgets"
	"github.com/raedahgroup/godcr-gio/widgets/editor"
)

type (
	passwordTab struct {
		passwordInput             *editor.Input
		confirmPasswordInput      *editor.Input
		passwordStrengthIndicator *widgets.ProgressBar
		passwordStrength          float64
	}

	pinTab struct {
		pinInput        *editor.Input
		confirmPinInput *editor.Input
		//	pinStrength     int
	}

	PinAndPasswordWidget struct {
		tabContainer *widgets.TabContainer
		passwordTab  *passwordTab
		pinTab       *pinTab

		currentTab string
		IsCreating bool

		createButton *widgets.Button
		cancelLabel  *widgets.ClickableLabel

		cancelFunc func()
		createFunc func()
	}
)

func NewPinAndPasswordWidget(cancelFunc, createFunc func()) *PinAndPasswordWidget {
	return &PinAndPasswordWidget{
		currentTab:   "password",
		tabContainer: widgets.NewTabContainer().AddTab("Password").AddTab("PIN"),
		createButton: widgets.NewButton("Create", nil),
		cancelLabel:  widgets.NewClickableLabel("Cancel").SetSize(4).SetWeight(text.Bold).SetColor(helper.DecredLightBlueColor),
		cancelFunc:   cancelFunc,
		createFunc:   createFunc,
		passwordTab: &passwordTab{
			passwordInput:             editor.NewInput("Spending Password").SetMask("*"),
			confirmPasswordInput:      editor.NewInput("Confirm Spending Password").SetMask("*"),
			passwordStrength:          0,
			passwordStrengthIndicator: widgets.NewProgressBar().SetHeight(6),
		},
		pinTab: &pinTab{
			pinInput:        editor.NewInput("Pin").SetMask("*").Numeric(),
			confirmPinInput: editor.NewInput("Confirm Pin").SetMask("*").Numeric(),
		},
	}
}

func (p *PinAndPasswordWidget) Reset() {
	p.passwordTab.passwordInput.Clear()
	p.passwordTab.confirmPasswordInput.Clear()
}

func (p *PinAndPasswordWidget) Value() string {
	if p.currentTab == "password" {
		return p.passwordTab.passwordInput.Text()
	}

	return p.pinTab.pinInput.Text()
}

func (p *PinAndPasswordWidget) Render(ctx *layout.Context) {
	layout.Stack{Alignment: layout.NW}.Layout(ctx,
		layout.Expanded(func() {
			widgets.NewLabel("Create a Spending Password", 5).
				SetWeight(text.Bold).
				Draw(ctx)
		}),
		layout.Stacked(func() {
			inset := layout.Inset{
				Top: unit.Dp(25),
			}
			inset.Layout(ctx, func() {
				p.tabContainer.Draw(ctx, p.passwordRenderFunc, p.pinRenderFunc)
			})
		}),
	)
}

func (p *PinAndPasswordWidget) passwordRenderFunc(ctx *layout.Context) {
	p.currentTab = "password"

	bothPasswordsMatch := false
	if p.passwordTab.confirmPasswordInput.Len() > 0 {
		if p.passwordTab.confirmPasswordInput.Text() == p.passwordTab.passwordInput.Text() {
			bothPasswordsMatch = true
		} else {
			bothPasswordsMatch = false
		}
	}

	// password section
	go func() {
		p.passwordTab.passwordStrength = (dcrlibwallet.ShannonEntropy(p.passwordTab.passwordInput.Text()) / 4) * 100
	}()

	p.passwordTab.passwordInput.Draw(ctx)

	// password strength section
	inset := layout.Inset{
		Top: unit.Dp(65),
	}
	inset.Layout(ctx, func() {
		var col color.RGBA
		if p.passwordTab.passwordStrength > 70 {
			col = helper.DecredGreenColor
		} else {
			col = helper.DecredOrangeColor
		}
		p.passwordTab.passwordStrengthIndicator.SetProgressColor(col).Draw(ctx, &p.passwordTab.passwordStrength)
	})

	// confirm password section
	inset = layout.Inset{
		Top: unit.Dp(85),
	}
	inset.Layout(ctx, func() {
		borderColor := helper.GrayColor
		focusBorderColor := helper.DecredLightBlueColor

		if !bothPasswordsMatch && p.passwordTab.confirmPasswordInput.Len() > 0 {
			borderColor = helper.DangerColor
			focusBorderColor = helper.DangerColor
		}
		p.passwordTab.confirmPasswordInput.SetBorderColor(borderColor).SetFocusedBorderColor(focusBorderColor).Draw(ctx)
	})

	// error text section
	inset = layout.Inset{
		Top: unit.Dp(145),
	}
	inset.Layout(ctx, func() {
		if !bothPasswordsMatch && p.passwordTab.confirmPasswordInput.Len() > 0 {
			widgets.NewLabel("Both passwords do not match").SetColor(helper.DangerColor).Draw(ctx)
		}
	})

	inset.Layout(ctx, func() {
		ctx.Constraints.Width.Min = ctx.Constraints.Width.Max
		layout.Stack{Alignment: layout.NE}.Layout(ctx,
			layout.Stacked(func() {
				layout.Flex{Axis: layout.Horizontal}.Layout(ctx,
					layout.Rigid(func() {
						inset := layout.Inset{
							Right: unit.Dp(10),
							Top:   unit.Dp(10),
						}
						inset.Layout(ctx, func() {
							var col color.RGBA
							if p.IsCreating {
								col = helper.GrayColor
							} else {
								col = helper.DecredLightBlueColor
							}

							p.cancelLabel.
								SetColor(col).
								Draw(ctx, func() {
									if !p.IsCreating {
										p.cancelFunc()
									}
								})
						})
					}),
					layout.Rigid(func() {
						createButton := p.createButton

						var txt string
						var bgCol color.RGBA
						if p.IsCreating {
							bgCol = helper.GrayColor
							txt = "Creating..."
						} else {
							txt = "Create"
							if bothPasswordsMatch && p.passwordTab.confirmPasswordInput.Len() > 0 {
								bgCol = helper.DecredLightBlueColor
							} else {
								bgCol = helper.GrayColor
							}
						}

						createButton.SetBackgroundColor(bgCol).SetText(txt).Draw(ctx, func() {
							if !p.IsCreating {
								p.createFunc()
							}
						})
					}),
				)
			}),
		)
	})

}

func (p *PinAndPasswordWidget) pinRenderFunc(ctx *layout.Context) {
	p.currentTab = "pin"
}
