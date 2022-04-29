package app

type Navigator interface {
	ChangePage(Page, bool) // Push / Pop / Replace
	ShowModal(Modal)
	DismissModal(Modal)
}

// ChangePage displays the provided page on the window and optionally adds
// the current page to the backstack. This automatically refreshes the display,
// callers should not re-refresh the display.
// Satisfies the Navigator interface.
func (app *App) ChangePage(page Page, keepBackStack bool) {
	if app.currentPage != nil && keepBackStack {
		app.currentPage.OnNavigatedFrom()
		app.pageBackStack = append(app.pageBackStack, app.currentPage)
	}

	app.currentPage = page
	app.currentPage.OnNavigatedTo()
	app.window.Invalidate()
}

// TODO: showModal should refresh display, callers shouldn't.
func (app *App) ShowModal(modal Modal) {
	modal.OnResume() // setup display data
	app.modalMutex.Lock()
	app.modals = append(app.modals, modal)
	app.modalMutex.Unlock()
}

func (app *App) DismissModal(modal Modal) {
	app.modalMutex.Lock()
	defer app.modalMutex.Unlock()
	for i, m := range app.modals {
		if m.ModalID() == modal.ModalID() {
			modal.OnDismiss() // do garbage collection in modal
			app.modals = append(app.modals[:i], app.modals[i+1:]...)
			app.window.Invalidate()
			return
		}
	}
}
