package app

import (
	"sync"
)

// PageStack is a stack of pages that handles page data initialization and
// destruction when pages are added to/removed from the top of the stack.
// NOTE: This stack does not maintain duplicate instances of the same page.
type PageStack struct {
	name  string
	mtx   sync.Mutex
	pages []Page
}

func NewPageStack(name string) *PageStack {
	return &PageStack{
		name: name,
	}
}

// Top returns the page that is at the top of the stack. Returns nil if the
// stack is empty.
func (pageStack *PageStack) Top() Page {
	pageStack.mtx.Lock()
	defer pageStack.mtx.Unlock()

	if l := len(pageStack.pages); l > 0 {
		return pageStack.pages[l-1]
	}
	return nil
}

// Push pushes the specified page to the top of the stack, removing all other
// instances of the same page from the stack. An about-to-display signal is sent
// to the new page via newPage.OnNavigatedTo() while page.OnNavigatedFrom() is
// called on the current page to signal that the current page is no longer the
// displayed page.
func (pageStack *PageStack) Push(newPage Page, navigator PageNavigator) {
	pageStack.mtx.Lock()
	defer pageStack.mtx.Unlock()

	if l := len(pageStack.pages); l > 0 {
		currentPage := pageStack.pages[l-1]
		if currentPage.ID() != newPage.ID() {
			currentPage.OnNavigatedFrom() // if currentPage is another instance of newPage, it'll be closed below.
		}
	}

	// Close any previous instance of this type. Use the Closed() method if
	// implemented, to signal that the instance will never be re-displayed.
	for i, existingPage := range pageStack.pages {
		if existingPage.ID() == newPage.ID() {
			existingPage.OnNavigatedFrom()
			if closablePage, ok := existingPage.(Closable); ok {
				closablePage.OnClosed()
			}
			pageStack.pages = append(pageStack.pages[:i], pageStack.pages[i+1:]...) // remove duplicate
			break
		}
	}

	pageStack.pages = append(pageStack.pages, newPage)
	newPage.OnAttachedToNavigator(navigator)
	newPage.OnNavigatedTo()
}

// Pop removes the page at the top of the stack and gets the next page ready for
// display. The OnNavigatedFrom() and if supported, the OnClosed() methods of
// the page to be removed are called to signal that the page is removed from
// display and will never be re-displayed. An about-to-display signal is sent to
// the page that will be displayed next via the page.OnNavigatedTo() method.
func (pageStack *PageStack) Pop() bool {
	pageStack.mtx.Lock()
	defer pageStack.mtx.Unlock()

	l := len(pageStack.pages)
	if l == 0 {
		return false
	}

	pageToPop := pageStack.pages[l-1]
	pageToPop.OnNavigatedFrom()
	if closeablePage, ok := pageToPop.(Closable); ok {
		closeablePage.OnClosed()
	}

	pageStack.pages = pageStack.pages[:l-1]
	if l > 1 {
		pageStack.pages[l-2].OnNavigatedTo() // get previous page ready for display
	}
	return true
}

// PopAfter removes all pages from the top of the stack until (and excluding) a
// specific page. The matcher parameter should return true for the page that
// should be excluded. If the matcher never matches a page to exclude, no page
// will be popped. If any page is popped, the page's OnNavigatedFrom() and if
// supported, the OnClosed() methods will be called to signal that the page has
// been removed from the display and will never be re-displayed. The page to be
// displayed will receive an about-to-display signal via the OnNavigatedTo()
// method.
func (pageStack *PageStack) PopAfter(matcher func(Page) bool) bool {
	pageStack.mtx.Lock()
	defer pageStack.mtx.Unlock()

	retainPageIndex := -1
	for i := len(pageStack.pages) - 1; i >= 0; i-- {
		if matcher(pageStack.pages[i]) {
			retainPageIndex = i
			break
		}
	}
	if retainPageIndex == -1 {
		return false
	}

	popped := pageStack.pages[retainPageIndex+1:] // pop pages after the retainPageIndex
	for _, poppedPage := range popped {
		poppedPage.OnNavigatedFrom()
		if closeablePage, ok := poppedPage.(Closable); ok {
			closeablePage.OnClosed()
		}
	}

	pageStack.pages = pageStack.pages[:retainPageIndex+1] // keep pages from index 0 up till retainPageIndex
	pageStack.pages[retainPageIndex].OnNavigatedTo()
	return true
}

// Reset pops all pages in the stack and creates a new stack with the specified
// pages as root. Each popped page's OnNavigatedFrom() and if supported, the
// OnClosed() methods will be called to signal that the page has been removed
// from the display and will never be re-displayed. If there are new pages to
// display, the top page is readied for display via the its OnNavigatedTo()
// method.
func (pageStack *PageStack) Reset(newPages ...Page) {
	pageStack.mtx.Lock()
	defer pageStack.mtx.Unlock()

	// Close all the pages in the current stack before resetting.
	for _, existingPage := range pageStack.pages {
		existingPage.OnNavigatedFrom()
		if closeablePage, ok := existingPage.(Closable); ok {
			closeablePage.OnClosed()
		}
	}

	pageStack.pages = newPages
	if l := len(newPages); l > 0 {
		pageStack.pages[l-1].OnNavigatedTo()
	}
}

func (pageStack *PageStack) pagesAfter(stopPageID *string) (pages []Page) {
	pageStack.mtx.Lock()
	defer pageStack.mtx.Unlock()

	for i := len(pageStack.pages) - 1; i >= 0; i-- {
		page := pageStack.pages[i]
		if stopPageID != nil && page.ID() == *stopPageID {
			break
		}
		pages = append(pages, page)
	}

	return pages
}
