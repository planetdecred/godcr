package materialplus

// type editor struct {
// 	widget   *widget.Editor
// 	material material.Editor
// 	line     *Line
// }

// // Password represents a form for collecting and confirming user password
// type Password struct {
// 	theme *Theme

// 	titleLabel material.Label

// 	spendingEditor *editor
// 	confirmEditor  *editor

// 	cancelLabelColor            color.RGBA
// 	createButtonBackgroundColor color.RGBA

// 	cancelButton         *widget.Button
// 	createButton         *widget.Button
// 	cancelButtonMaterial material.Button
// 	createButtonMaterial material.Button

// 	errorLabel material.Label
// }

// // Password initializes and returns an instance of Password
// func (t *Theme) Password() *Password {
// 	cancelButtonMaterial := t.Button("Cancel")
// 	cancelButtonMaterial.Background = color.RGBA{}
// 	cancelButtonMaterial.Color = t.Color.Primary

// 	errorLabel := t.Body2("")
// 	errorLabel.Color = t.Danger

// 	spendingEditor := &editor{
// 		widget:   new(widget.Editor),
// 		material: t.Editor("Spending password"),
// 		line:     t.Line(),
// 	}
// 	confirmEditor := &editor{
// 		widget:   new(widget.Editor),
// 		material: t.Editor("Confirm spending password"),
// 		line:     t.Line(),
// 	}

// 	confirmEditor.line.Color = t.Color.Primary
// 	spendingEditor.line.Color = t.Color.Primary

// 	p := &Password{
// 		theme: t,

// 		titleLabel: t.H5("Create spending password"),
// 		errorLabel: errorLabel,

// 		confirmEditor:        confirmEditor,
// 		spendingEditor:       spendingEditor,
// 		cancelButton:         new(widget.Button),
// 		createButton:         new(widget.Button),
// 		cancelButtonMaterial: cancelButtonMaterial,
// 		createButtonMaterial: t.Button("Create"),
// 	}

// 	return p
// }

// func (p *Password) updateColors() {
// 	p.cancelLabelColor = p.theme.Disabled
// 	p.createButtonBackgroundColor = p.theme.Disabled

// 	if p.bothPasswordsMatch() && p.confirmEditor.widget.Len() > 0 {
// 		p.createButtonBackgroundColor = p.theme.Color.Primary
// 	}
// }

// func (p *Password) processButtonClicks(gtx *layout.Context, confirm func(string), cancel func()) {
// 	for p.createButton.Clicked(gtx) {
// 		if p.validate() {
// 			confirm(p.spendingEditor.widget.Text())
// 		}
// 	}

// 	for p.cancelButton.Clicked(gtx) {
// 		p.Reset()
// 		cancel()
// 	}
// }

// // Draw renders the widget to screen. The confirm function passed by the calling page is called when the confirm button
// // is clicked, and the form passes validation. The entered password is passed as an argument to the confirm func.
// // The cancel func is called when the cancel button is clicked
// func (p *Password) Draw(gtx *layout.Context, confirm func(string), cancel func()) {
// 	p.processButtonClicks(gtx, confirm, cancel)
// 	p.updateColors()
// 	p.validate()

// 	widgets := []func(){
// 		func() {
// 			p.titleLabel.Layout(gtx)
// 		},
// 		func() {
// 			inset := layout.Inset{
// 				Top: unit.Dp(20),
// 			}
// 			inset.Layout(gtx, func() {
// 				p.spendingEditor.layout(gtx)
// 			})
// 		},
// 		func() {
// 			p.confirmEditor.layout(gtx)
// 		},
// 		func() {
// 			layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
// 				layout.Flexed(2, func() {
// 					p.errorLabel.Layout(gtx)
// 				}),
// 				layout.Rigid(func() {
// 					layout.UniformInset(unit.Dp(10)).Layout(gtx, func() {
// 						p.cancelButtonMaterial.Color = p.cancelLabelColor
// 						p.cancelButtonMaterial.Layout(gtx, p.cancelButton)
// 					})
// 				}),
// 				layout.Rigid(func() {
// 					gtx.Constraints.Width.Min = 70
// 					layout.UniformInset(unit.Dp(10)).Layout(gtx, func() {
// 						p.createButtonMaterial.Background = p.createButtonBackgroundColor
// 						p.createButtonMaterial.Layout(gtx, p.createButton)
// 					})
// 				}),
// 			)
// 		},
// 	}

// 	list := layout.List{Axis: layout.Vertical}
// 	list.Layout(gtx, len(widgets), func(i int) {
// 		layout.UniformInset(unit.Dp(0)).Layout(gtx, widgets[i])
// 	})
// }

// func (p *Password) validate() bool {
// 	p.errorLabel.Text = ""

// 	if p.spendingEditor.widget.Len() == 0 || p.confirmEditor.widget.Len() == 0 {
// 		return false
// 	}

// 	if p.spendingEditor.widget.Len() > 0 && p.confirmEditor.widget.Len() > 0 && !p.bothPasswordsMatch() {
// 		p.errorLabel.Text = "Both passwords do not match"
// 		return false
// 	}

// 	return true
// }

// func (p *Password) bothPasswordsMatch() bool {
// 	return p.confirmEditor.widget.Text() == p.spendingEditor.widget.Text()
// }

// // Reset empties the contents of the password form
// func (p *Password) Reset() {
// 	p.spendingEditor.widget.SetText("")
// 	p.confirmEditor.widget.SetText("")
// }

// func (e *editor) layout(gtx *layout.Context) {
// 	e.material.Layout(gtx, e.widget)
// 	inset := layout.Inset{
// 		Top: unit.Dp(float32(gtx.Constraints.Height.Min)),
// 	}
// 	inset.Layout(gtx, func() {
// 		e.line.Draw(gtx)
// 	})
// }
