package app

// GenericPageModal implements the ID() and OnAttachedToNavigator() methods
// required by most pages and modals. It also defines ParentNavigator() and
// ParentWindow() helper methods, to enable pages access the Navigator that
// displayed the page and the root WindowNavigator.
// Actual pages and modals may embed this struct and implement other methods
// as necessary.
type GenericPageModal struct {
	id        string
	parentNav PageNavigator
}

// NewGenericPageModal returns an instance of a GenericPageModal.
func NewGenericPageModal(id string) *GenericPageModal {
	return &GenericPageModal{
		id: id,
	}
}

// ID is a unique string that identifies this page or modal and may be used to
// differentiate this page or modal from other pages or modals.
// Part of the Page and Modal interfaces.
func (pageModal *GenericPageModal) ID() string {
	return pageModal.id
}

// OnAttachedToNavigator is called when navigation occurs; i.e. when this page
// or modal is pushed into the window's display. The navigator parameter is the
// PageNavigator or WindowNavigator object that is used to display this page or
// modal. OnAttachedToNavigator is called just before OnResume (for modals) and
// OnNavigatedTo (for pages).
// Part of the Page and Modal interfaces.
func (pageModal *GenericPageModal) OnAttachedToNavigator(parentNav PageNavigator) {
	pageModal.parentNav = parentNav
}

// ParentNavigator is a helper method that returns the Navigator that pushed
// this content into display, which may be the WindowNavigator or any other page
// that implements the PageNavigator interface (e.g. a MasterPage). For modals,
// this is always the WindowNavigator.
func (pageModal *GenericPageModal) ParentNavigator() PageNavigator {
	return pageModal.parentNav
}

// ParentWindow is a helper method that returns the Navigator that displayed
// this page or modal if it is a WindowNavigator, otherwise it recursively
// checks the parent navigators to find and return a WindowNavigator.
func (pageModal *GenericPageModal) ParentWindow() WindowNavigator {
	parentNav := pageModal.ParentNavigator()
	for {
		if parentNav == nil {
			return nil
		}
		if windowNav, isWindowNav := parentNav.(WindowNavigator); isWindowNav {
			return windowNav
		}
		if navigatedPageModal, ok := parentNav.(interface{ ParentNavigator() PageNavigator }); ok {
			parentNav = navigatedPageModal.ParentNavigator()
		}
	}
}
