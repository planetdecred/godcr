package wallets

import (
	"gioui.org/io/clipboard"
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/app"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/modal"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

const SignMessagePageID = "SignMessage"

type SignMessagePage struct {
	*app.App
	PopFragment func() // TODO: Will crash.

	container layout.List
	wallet    *dcrlibwallet.Wallet

	isSigningMessage bool
	addressIsValid   bool
	messageIsValid   bool
	isEnabled        bool

	titleLabel, errorLabel, signedMessageLabel decredmaterial.Label
	addressEditor, messageEditor               decredmaterial.Editor
	clearButton, signButton, copyButton        decredmaterial.Button
	copySignature                              *decredmaterial.Clickable
	copyIcon                                   *decredmaterial.Image
	gtx                                        *layout.Context

	backButton decredmaterial.IconButton
	infoButton decredmaterial.IconButton
}

func NewSignMessagePage(app *app.App, wallet *dcrlibwallet.Wallet) *SignMessagePage {
	addressEditor := app.Theme.Editor(new(widget.Editor), "Address")
	addressEditor.Editor.SingleLine, addressEditor.Editor.Submit = true, true
	messageEditor := app.Theme.Editor(new(widget.Editor), "Message")
	messageEditor.Editor.SingleLine, messageEditor.Editor.Submit = true, true

	clearButton := app.Theme.OutlineButton("Clear all")
	signButton := app.Theme.Button("Sign message")
	clearButton.Font.Weight, signButton.Font.Weight = text.Medium, text.Medium
	signButton.SetEnabled(false)
	clearButton.SetEnabled(false)

	errorLabel := app.Theme.Caption("")
	errorLabel.Color = app.Theme.Color.Danger
	copyIcon := app.Theme.Icons.CopyIcon

	pg := &SignMessagePage{
		App:    app,
		wallet: wallet,
		container: layout.List{
			Axis: layout.Vertical,
		},

		titleLabel:         app.Theme.H5("Sign Message"),
		signedMessageLabel: app.Theme.Body1(""),
		errorLabel:         errorLabel,
		addressEditor:      addressEditor,
		messageEditor:      messageEditor,
		clearButton:        clearButton,
		signButton:         signButton,
		copyButton:         app.Theme.Button("Copy"),
		copySignature:      app.Theme.NewClickable(false),
		copyIcon:           copyIcon,
	}

	pg.signedMessageLabel.Color = app.Theme.Color.GrayText2
	pg.backButton, pg.infoButton = components.SubpageHeaderButtons(app.Theme)

	return pg
}

// ID is a unique string that identifies the page and may be used
// to differentiate this page from other pages.
// Part of the load.Page interface.
func (pg *SignMessagePage) ID() string {
	return SignMessagePageID
}

// OnNavigatedTo is called when the page is about to be displayed and
// may be used to initialize page features that are only relevant when
// the page is displayed.
// Part of the load.Page interface.
func (pg *SignMessagePage) OnNavigatedTo() {
	pg.addressEditor.Editor.Focus()
}

// Layout draws the page UI components into the provided layout context
// to be eventually drawn on screen.
// Part of the load.Page interface.
func (pg *SignMessagePage) Layout(gtx layout.Context) layout.Dimensions {
	body := func(gtx C) D {
		sp := components.SubPage{
			App:        pg.App,
			Title:      "Sign message",
			WalletName: pg.wallet.Name,
			BackButton: pg.backButton,
			InfoButton: pg.infoButton,
			Back: func() {
				pg.PopFragment()
			},
			Body: func(gtx layout.Context) layout.Dimensions {
				return pg.Theme.Card().Layout(gtx, func(gtx C) D {
					return layout.UniformInset(values.MarginPadding15).Layout(gtx, func(gtx C) D {
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
							layout.Rigid(pg.description()),
							layout.Rigid(pg.editors(pg.addressEditor)),
							layout.Rigid(pg.editors(pg.messageEditor)),
							layout.Rigid(pg.drawButtonsRow()),
							layout.Rigid(pg.drawResult()),
						)
					})
				})
			},
			InfoTemplate: modal.SignMessageInfoTemplate,
		}
		return sp.Layout(gtx)
	}

	return components.UniformPadding(gtx, body)
}

func (pg *SignMessagePage) description() layout.Widget {
	return func(gtx C) D {
		desc := pg.Theme.Caption("Enter an address and message to sign:")
		desc.Color = pg.Theme.Color.GrayText2
		return layout.Inset{Bottom: values.MarginPadding20}.Layout(gtx, desc.Layout)
	}
}

func (pg *SignMessagePage) editors(editor decredmaterial.Editor) layout.Widget {
	return func(gtx C) D {
		return layout.Inset{Bottom: values.MarginPadding15}.Layout(gtx, editor.Layout)
	}
}

func (pg *SignMessagePage) drawButtonsRow() layout.Widget {
	return func(gtx C) D {
		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Flexed(1, func(gtx C) D {
				return layout.E.Layout(gtx, func(gtx C) D {
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							inset := layout.Inset{
								Right: values.MarginPadding5,
							}
							return inset.Layout(gtx, pg.clearButton.Layout)
						}),
						layout.Rigid(pg.signButton.Layout),
					)
				})
			}),
		)
	}
}

func (pg *SignMessagePage) drawResult() layout.Widget {
	return func(gtx C) D {
		if !components.StringNotEmpty(pg.signedMessageLabel.Text) {
			return layout.Dimensions{}
		}
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				m := values.MarginPadding30
				return layout.Inset{Top: m, Bottom: m}.Layout(gtx, pg.Theme.Separator().Layout)
			}),
			layout.Rigid(func(gtx C) D {
				return layout.Stack{}.Layout(gtx,
					layout.Stacked(func(gtx C) D {
						border := widget.Border{Color: pg.Theme.Color.Gray4, CornerRadius: values.MarginPadding10, Width: values.MarginPadding2}
						wrapper := pg.Theme.Card()
						wrapper.Color = pg.Theme.Color.Gray4
						return border.Layout(gtx, func(gtx C) D {
							return wrapper.Layout(gtx, func(gtx C) D {
								return layout.UniformInset(values.MarginPadding10).Layout(gtx, func(gtx C) D {
									return layout.Flex{}.Layout(gtx,
										layout.Flexed(0.9, pg.signedMessageLabel.Layout),
										layout.Flexed(0.1, func(gtx C) D {
											return layout.E.Layout(gtx, func(gtx C) D {
												return layout.Inset{Top: values.MarginPadding7}.Layout(gtx, func(gtx C) D {
													if pg.copySignature.Clicked() {
														clipboard.WriteOp{Text: pg.signedMessageLabel.Text}.Add(gtx.Ops)
														pg.Toast.Notify("Signature copied")
													}
													return pg.copySignature.Layout(gtx, pg.copyIcon.Layout24dp)
												})
											})
										}),
									)
								})
							})
						})
					}),
					layout.Stacked(func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{
							Top:  values.MarginPaddingMinus10,
							Left: values.MarginPadding10,
						}.Layout(gtx, func(gtx C) D {
							return pg.Theme.Card().Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								label := pg.Theme.Body1("Signature")
								label.Color = pg.Theme.Color.GrayText2
								return label.Layout(gtx)
							})
						})
					}),
				)
			}),
		)
	}
}

func (pg *SignMessagePage) updateButtonColors() {
	pg.isEnabled = false
	if components.StringNotEmpty(pg.addressEditor.Editor.Text()) ||
		components.StringNotEmpty(pg.messageEditor.Editor.Text()) {
		pg.clearButton.SetEnabled(true)
	} else {
		pg.clearButton.SetEnabled(false)
	}

	if !pg.isSigningMessage && pg.messageIsValid && pg.addressIsValid {
		pg.isEnabled = true
	}

	pg.signButton.SetEnabled(pg.isEnabled)
}

// HandleUserInteractions is called just before Layout() to determine
// if any user interaction recently occurred on the page and may be
// used to update the page's UI components shortly before they are
// displayed.
// Part of the load.Page interface.
func (pg *SignMessagePage) HandleUserInteractions() {
	pg.updateButtonColors()

	isSubmit, isChanged := decredmaterial.HandleEditorEvents(pg.addressEditor.Editor, pg.messageEditor.Editor)
	if isChanged {
		if pg.addressEditor.Editor.Focused() {
			pg.validateAddress()
		}

		if pg.messageEditor.Editor.Focused() {
			pg.validateMessage()
		}
	}

	for pg.clearButton.Clicked() {
		pg.clearForm()
	}

	if (pg.signButton.Clicked() || isSubmit) && pg.isEnabled {
		if !pg.isSigningMessage && pg.validate() {
			address := pg.addressEditor.Editor.Text()
			message := pg.messageEditor.Editor.Text()

			modal.NewPasswordModal(pg.Theme, pg.App).
				Title("Confirm to sign").
				NegativeButton("Cancel", func() {}).
				PositiveButton("Confirm", func(password string, pm *modal.PasswordModal) bool {
					go func() {
						sig, err := pg.wallet.SignMessage([]byte(password), address, message)
						if err != nil {
							pm.SetError(err.Error())
							pm.SetLoading(false)
							return
						}

						pm.Dismiss()
						pg.signedMessageLabel.Text = dcrlibwallet.EncodeBase64(sig)

					}()
					return false
				}).Show()
		}
	}
}

// HandleKeyEvent is called when a key is pressed on the current window.
// Satisfies the app.KeyEventHandler interface for receiving key events.
func (pg *SignMessagePage) HandleKeyEvent(evt *key.Event) {
	// Switch editors when tab key is pressed.
	decredmaterial.SwitchEditors(evt, pg.addressEditor.Editor, pg.messageEditor.Editor)
}

func (pg *SignMessagePage) validate() bool {
	if !pg.validateAddress() || !pg.validateMessage() {
		return false
	}
	return true
}

func (pg *SignMessagePage) validateAddress() bool {
	address := pg.addressEditor.Editor.Text()
	pg.addressEditor.SetError("")

	var valid bool
	var errorMessage string

	switch {
	case !components.StringNotEmpty(address):
		errorMessage = "Please enter a valid address"
	case !pg.MultiWallet().IsAddressValid(address):
		errorMessage = "Invalid address"
	case !pg.wallet.HaveAddress(address):
		errorMessage = "Address not owned by any wallet"
	default:
		valid = true
	}
	if !valid {
		pg.addressEditor.SetError(errorMessage)
	}

	pg.addressIsValid = valid
	return valid
}

func (pg *SignMessagePage) validateMessage() bool {
	message := pg.messageEditor.Editor.Text()
	pg.messageEditor.SetError("")

	if !components.StringNotEmpty(message) {
		pg.messageEditor.SetError("Please enter a valid message to sign")
		pg.messageIsValid = false
		return false
	}
	pg.messageIsValid = true
	return true
}

func (pg *SignMessagePage) clearForm() {
	pg.addressEditor.Editor.SetText("")
	pg.messageEditor.Editor.SetText("")
	pg.addressEditor.SetError("")
	pg.messageEditor.SetError("")
	pg.signedMessageLabel.Text = ""
	pg.errorLabel.Text = ""
}

// OnNavigatedFrom is called when the page is about to be removed from
// the displayed window. This method should ideally be used to disable
// features that are irrelevant when the page is NOT displayed.
// NOTE: The page may be re-displayed on the app's window, in which case
// OnNavigatedTo() will be called again. This method should not destroy UI
// components unless they'll be recreated in the OnNavigatedTo() method.
// Part of the load.Page interface.
func (pg *SignMessagePage) OnNavigatedFrom() {}
