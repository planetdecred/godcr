package ui

func (win *Window) Overview() {
	body := func() {
		win.outputs.notImplemented.Layout(win.gtx)
	}
	win.Page(body)
}
