package app

// MasterPage is a page that can display subpages. It is an extension of the
// GenericPageModal which provides access to the Window or PageNavigator that
// was used to display the MasterPage. The ParentNavigator of a MasterPage is
// typically set when the MasterPage is pushed into the display window by a
// WindowNavigator or a PageNavigator.
// MasterPage must be embedded by pages that want to display subpages. Those
// pages must satisfy the other methods of the Page interface that are not
// already satisifed by MasterPage.
type MasterPage struct {
	*GenericPageModal
	subPages *PageStack
}

// NewMasterPage returns an instance of MasterPage.
func NewMasterPage(id string) *MasterPage {
	return &MasterPage{
		GenericPageModal: NewGenericPageModal(id),
		subPages:         NewPageStack(id),
	}
}

// CurrentPage returns the page that is at the top of the stack. Returns nil if
// the stack is empty.
// Part of the PageNavigator interface.
func (masterPage *MasterPage) CurrentPage() Page {
	return masterPage.subPages.Top()
}

// CurrentPageID returns the ID of the current page or an empty string if no
// page is displayed.
// Part of the PageNavigator interface.
func (masterPage *MasterPage) CurrentPageID() string {
	if currentPage := masterPage.CurrentPage(); currentPage != nil {
		return currentPage.ID()
	}
	return ""
}

// Display causes the specified page to be displayed on the parent window or
// page. All other instances of this same page will be closed and removed
// from the backstack.
// Part of the PageNavigator interface.
func (masterPage *MasterPage) Display(newPage Page) {
	newPage.OnAttachedToNavigator(masterPage)
	masterPage.subPages.Push(newPage)
	masterPage.ParentWindow().Reload()
}

// CloseCurrentPage dismisses the page at the top of the stack and gets the next
// page ready for display.
// Part of the PageNavigator interface.
func (masterPage *MasterPage) CloseCurrentPage() {
	popped := masterPage.subPages.Pop()
	if popped {
		masterPage.ParentWindow().Reload()
	}
}

// ClosePagesAfter dismisses all pages from the top of the stack until (and
// excluding) the page with the specified ID. If no page is found with the
// provided ID, no page will be popped. The page with the specified ID will be
// displayed after the other pages are popped.
// Part of the PageNavigator interface.
func (masterPage *MasterPage) ClosePagesAfter(keepPageID string) {
	popped := masterPage.subPages.PopAfter(func(page Page) bool {
		return page.ID() == keepPageID
	})
	if popped {
		masterPage.ParentWindow().Reload()
	}
}

// ClearStackAndDisplay dismisses all pages in the stack and displays the
// specified page.
// Part of the PageNavigator interface.
func (masterPage *MasterPage) ClearStackAndDisplay(newPage Page) {
	newPage.OnAttachedToNavigator(masterPage)
	masterPage.subPages.Reset(newPage)
	masterPage.ParentWindow().Reload()
}

// CloseAllPages dismisses all pages in the stack.
// Part of the PageNavigator interface.
func (masterPage *MasterPage) CloseAllPages() {
	masterPage.subPages.Reset()
	masterPage.ParentWindow().Reload()
}
