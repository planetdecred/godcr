package ui

func (win *Window) Receive() {
	body := func() {
		win.outputs.notImplemented.Layout(win.gtx)
	}
	win.Page(body)
}
