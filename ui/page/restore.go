package page

import (
	"strconv"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	. "github.com/raedahgroup/godcr-gio/ui"
	"github.com/raedahgroup/godcr-gio/ui/helper"
	"github.com/raedahgroup/godcr-gio/ui/themes/materialplus"
	"github.com/raedahgroup/godcr-gio/ui/units"
	"github.com/raedahgroup/godcr-gio/wallet"
)

// RestoreID is the id of the restore page.
const RestoreID = "restore"

// Restore represents the restore page of the app.
type Restore struct {
	container layout.Flex
	heading   material.Label
	theme     *materialplus.Theme

	restoreBtn material.Button
	restoreWdg *widget.Button

	backBtn material.Button
	backWdg *widget.Button

	inputGroupContainer layout.List
	inputGroupHeader    material.Label
	inputLabels         []material.Label
	inputs              []material.Editor
	editors             []*widget.Editor

	modal
	showModal bool
}

type modal struct {
	inputPassword         material.Editor
	editorPassword        *widget.Editor
	inputConfirmPassword  material.Editor
	editorConfirmPassword *widget.Editor
	inputSendingPin       material.Editor
	editorSendingPin      *widget.Editor
	cancelBtn             material.Button
	cancelWdg             *widget.Button
	submitBtn             material.Button
	submitWdg             *widget.Button
	confirmErrorMsg       material.Label
}

// Init initializes restore page with heading, 33 inputs and button
func (pg *Restore) Init(theme *materialplus.Theme, _ *wallet.Wallet, states map[string]interface{}) {
	pg.heading = theme.Label(units.Label, "Restore from seed phrase")
	pg.heading.Alignment = text.Middle
	pg.theme = theme
	pg.container.Axis = layout.Vertical

	pg.inputGroupHeader = theme.Label(unit.Dp(16), "Enter your seed phrase in the correct order")
	pg.inputGroupHeader.Alignment = text.Middle
	pg.inputGroupHeader.Color = GrayColor
	pg.inputGroupContainer.Axis = layout.Vertical

	pg.restoreBtn = theme.Button("Restore")
	pg.restoreBtn.Background = LightBlueColor
	pg.restoreWdg = new(widget.Button)

	pg.backBtn = theme.Button("Back")
	pg.backBtn.Background = LightGrayColor
	pg.backWdg = new(widget.Button)

	for i := 0; i <= 32; i++ {
		pg.inputs = append(pg.inputs, theme.Editor("Input word "+strconv.Itoa(i+1)+"..."))
		pg.inputs[i].Font.Size = unit.Sp(16)
		pg.editors = append(pg.editors, &widget.Editor{SingleLine: true})
		pg.inputLabels = append(pg.inputLabels, theme.Label(unit.Dp(13), "Word #"+strconv.Itoa(i+1)))
	}

	pg.showModal = false
	pg.initModal(theme)
}

// initModal initializes modal with inputs password, confirm password, PIN
// and buttons
func (pg *modal) initModal(theme *materialplus.Theme) {
	pg.confirmErrorMsg = theme.Label(unit.Dp(16), "")
	pg.confirmErrorMsg.Color = DangerColor

	pg.inputPassword = theme.Editor("Password...")
	pg.editorPassword = &widget.Editor{SingleLine: true}
	pg.inputConfirmPassword = theme.Editor("Confirm password...")
	pg.editorConfirmPassword = &widget.Editor{SingleLine: true}
	pg.inputSendingPin = theme.Editor("Input PIN...")
	pg.editorSendingPin = &widget.Editor{SingleLine: true}

	pg.cancelBtn = theme.Button("Cancel")
	pg.cancelBtn.Background = LightGrayColor
	pg.cancelWdg = new(widget.Button)
	pg.submitBtn = theme.Button("Create")
	pg.submitBtn.Background = LightBlueColor
	pg.submitWdg = new(widget.Button)
}

// Draw renders the page widgets
func (pg *Restore) Draw(gtx *layout.Context) interface{} {
	layout.UniformInset(units.FlexInset).Layout(gtx, func() {
		pg.container.Layout(gtx,
			layout.Rigid(func() {
				gtx.Constraints.Width.Min = gtx.Constraints.Width.Max
				pg.heading.Layout(gtx)
			}),
			layout.Rigid(func() {
				gtx.Constraints.Width.Min = gtx.Constraints.Width.Max
				pg.inputGroupHeader.Layout(gtx)
			}),
			layout.Flexed(1, func() {
				pg.inputGroupContainer.Layout(gtx, len(pg.inputs),
					layout.ListElement(func(i int) {
						layout.Flex{Axis: layout.Horizontal, Alignment: layout.Baseline}.Layout(gtx,
							layout.Rigid(func() {
								pg.inputLabels[i].Layout(gtx)
							}),
							layout.Flexed(1, func() {
								gtx.Constraints.Width.Min = gtx.Constraints.Width.Max
								layout.Inset{Left: unit.Dp(20), Bottom: unit.Dp(20)}.Layout(gtx, func() {
									pg.inputs[i].Layout(gtx, pg.editors[i])
								})
							}),
						)
					}),
				)
			}),
			layout.Rigid(func() {
				gtx.Constraints.Width.Min = gtx.Constraints.Width.Max
				layout.Inset{Top: unit.Dp(10)}.Layout(gtx, func() {
					layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Rigid(func() {
							pg.backBtn.Layout(gtx, pg.backWdg)
						}),
						layout.Rigid(func() {
							layout.Inset{Left: unit.Dp(10)}.Layout(gtx, func() {
								pg.restoreBtn.Layout(gtx, pg.restoreWdg)
							})
						}),
					)
				})
			}),
		)
	})

	pg.drawModal(gtx)

	if pg.backWdg.Clicked(gtx) {
		return EventNav{
			Current: RestoreID,
			Next:    LandingID,
		}
	}

	if pg.restoreWdg.Clicked(gtx) {
		if pg.validateWords() {
			pg.showModal = true
		}
	}

	if pg.cancelWdg.Clicked(gtx) {
		pg.showModal = false
	}

	if pg.submitWdg.Clicked(gtx) {
		if pg.validatePwdAndPIN() {
			return EventNav{
				Current: RestoreID,
				Next:    WalletsID,
			}
		}
	}

	return nil
}

func (pg *Restore) drawModal(gtx *layout.Context) {
	if !pg.showModal {
		return
	}
	widgets := []func(){
		func() {
			layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func() {
					pg.theme.H5("Create a spending password").Layout(gtx)
				}),
				layout.Rigid(func() {
					layout.Inset{Top: unit.Dp(10)}.Layout(gtx, func() {})
				}),
				layout.Rigid(func() {
					layout.Flex{Axis: layout.Vertical, Alignment: layout.Baseline}.Layout(gtx,
						layout.Rigid(func() {
							pg.theme.Label(unit.Dp(16), "Spending password").Layout(gtx)
						}),
						layout.Rigid(func() {
							layout.Inset{Top: unit.Dp(8)}.Layout(gtx, func() {})
						}),
						layout.Rigid(func() {
							gtx.Constraints.Width.Min = gtx.Constraints.Width.Max
							pg.inputPassword.Layout(gtx, pg.editorPassword)
						}),
					)
				}),
				layout.Rigid(func() {
					layout.Inset{Top: unit.Dp(25)}.Layout(gtx, func() {})
				}),
				layout.Rigid(func() {
					layout.Flex{Axis: layout.Vertical, Alignment: layout.Baseline}.Layout(gtx,
						layout.Rigid(func() {
							pg.theme.Label(unit.Dp(16), "Confirm spending password").Layout(gtx)
						}),
						layout.Rigid(func() {
							layout.Inset{Top: unit.Dp(8)}.Layout(gtx, func() {})
						}),
						layout.Rigid(func() {
							gtx.Constraints.Width.Min = gtx.Constraints.Width.Max
							pg.inputConfirmPassword.Layout(gtx, pg.editorConfirmPassword)
						}),
					)
				}),
			)
		},
		func() {
			layout.Inset{Top: unit.Dp(15)}.Layout(gtx, func() {})
		},
		func() {
			helper.PaintArea(gtx, LightGrayColor, gtx.Constraints.Width.Max, 1)
			layout.Inset{Top: unit.Dp(15)}.Layout(gtx, func() {})
		},
		func() {
			layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func() {
					pg.theme.H5("Create a spending PIN").Layout(gtx)
				}),
				layout.Rigid(func() {
					layout.Inset{Top: unit.Dp(10)}.Layout(gtx, func() {})
				}),
				layout.Rigid(func() {
					layout.Flex{Axis: layout.Vertical, Alignment: layout.Baseline}.Layout(gtx,
						layout.Rigid(func() {
							pg.theme.Label(unit.Dp(16), "Enter spending PIN").Layout(gtx)
						}),
						layout.Rigid(func() {
							layout.Inset{Top: unit.Dp(8)}.Layout(gtx, func() {})
						}),
						layout.Rigid(func() {
							gtx.Constraints.Width.Min = gtx.Constraints.Width.Max
							pg.inputSendingPin.Layout(gtx, pg.editorSendingPin)
						}),
					)
				}),
			)
		},
	}

	layout.Inset{Top: unit.Dp(0), Left: unit.Dp(0)}.Layout(gtx, func() {
		helper.Fill(gtx, LightGrayColor)
		layout.Inset{Top: unit.Dp(120), Left: unit.Dp(0)}.Layout(gtx, func() {
			helper.Fill(gtx, WhiteColor)
			layout.UniformInset(units.FlexInset).Layout(gtx, func() {
				layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Flexed(1, func() {
						(&layout.List{Axis: layout.Vertical}).Layout(gtx, len(widgets), func(i int) {
							layout.Inset{}.Layout(gtx, widgets[i])
						})
					}),
					layout.Rigid(func() {
						gtx.Constraints.Width.Min = gtx.Constraints.Width.Max
						layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
							layout.Rigid(func() {
								pg.cancelBtn.Layout(gtx, pg.cancelWdg)
							}),
							layout.Rigid(func() {
								layout.Inset{Left: unit.Dp(10)}.Layout(gtx, func() {
									pg.submitBtn.Layout(gtx, pg.submitWdg)
								})
							}),
							layout.Rigid(func() {
								layout.Inset{Left: unit.Dp(20)}.Layout(gtx, func() {
									pg.confirmErrorMsg.Layout(gtx)
								})
							}),
						)
					}),
				)
			})
		})
	})
}

func (pg *Restore) validateWords() bool {
	for i, editor := range pg.editors {
		txt := editor.Text()
		pg.inputLabels[i].Color = BlackColor
		if txt == "" {
			pg.inputLabels[i].Color = DangerColor
			return false
		}
	}
	return true
}

func (pg *modal) validatePwdAndPIN() bool {
	pg.confirmErrorMsg.Text = ""
	if pg.editorPassword.Text() == "" {
		pg.confirmErrorMsg.Text = "Please enter password"
		return false
	}
	if pg.editorConfirmPassword.Text() == "" {
		pg.confirmErrorMsg.Text = "Please enter confirm password"
		return false
	}
	if pg.editorConfirmPassword.Text() != pg.editorPassword.Text() {
		pg.confirmErrorMsg.Text = "Password and confirm password do not match"
		return false
	}
	if pg.editorSendingPin.Text() == "" {
		pg.confirmErrorMsg.Text = "Please enter spending PIN"
		return false
	}
	return true
}
