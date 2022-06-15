package app

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"testing"

	"gioui.org/layout"
)

type tPageMethod uint8

const (
	pageMethodOnAttachedToNavigator tPageMethod = iota
	pageMethodOnNavigatedTo
	pageMethodHandleUserInteractions
	pageMethodLayout
	pageMethodOnNavigatedFrom
	pageMethodOnClosed
)

func combinePageMethods(methods []tPageMethod) string {
	methodNames := make([]string, len(methods))
	for i := range methods {
		methodNames[i] = methods[i].String()
	}
	return strings.Join(methodNames, ",")
}

func (tpm tPageMethod) String() string {
	switch tpm {
	case pageMethodOnAttachedToNavigator:
		return "OnAttachedToNavigator()"
	case pageMethodOnNavigatedTo:
		return "OnNavigatedTo()"
	case pageMethodHandleUserInteractions:
		return "HandleUserInteractions()"
	case pageMethodLayout:
		return "Layout()"
	case pageMethodOnNavigatedFrom:
		return "OnNavigatedFrom()"
	case pageMethodOnClosed:
		return "OnClosed()"
	default:
		return fmt.Sprintf("UnknownMethod %d", tpm)
	}
}

type testPage struct {
	id               string
	log              func(format string, args ...interface{})
	mtx              sync.Mutex
	parentNav        PageNavigator
	methodsCallStack []tPageMethod
}

func newTestPage(id string, logFn func(format string, args ...interface{})) *testPage {
	return &testPage{
		id:  id,
		log: logFn,
	}
}

func (tPage *testPage) calledMethods() []tPageMethod {
	tPage.mtx.Lock()
	defer tPage.mtx.Unlock()
	methods := tPage.methodsCallStack
	tPage.methodsCallStack = nil
	tPage.log("cleared called methods on %s", tPage.id)
	return methods
}

func (tPage *testPage) calledMethod(tpm tPageMethod) {
	tPage.mtx.Lock()
	defer tPage.mtx.Unlock()
	tPage.methodsCallStack = append(tPage.methodsCallStack, tpm)
	tPage.log("recorded %s.%s", tPage.id, tpm)
}

func (tPage *testPage) ID() string {
	// No need to track calls to the ID() method, doesn't do anything special.
	return tPage.id
}

func (tPage *testPage) OnAttachedToNavigator(navigator PageNavigator) {
	tPage.calledMethod(pageMethodOnAttachedToNavigator)
	tPage.parentNav = navigator
}

func (tPage *testPage) OnNavigatedTo() {
	tPage.calledMethod(pageMethodOnNavigatedTo)
}

func (tPage *testPage) HandleUserInteractions() {
	tPage.calledMethod(pageMethodHandleUserInteractions)
}

func (tPage *testPage) Layout(layout.Context) layout.Dimensions {
	tPage.calledMethod(pageMethodLayout)
	return layout.Dimensions{}
}

func (tPage *testPage) OnNavigatedFrom() {
	tPage.calledMethod(pageMethodOnNavigatedFrom)
}

func (tPage *testPage) OnClosed() {
	tPage.calledMethod(pageMethodOnClosed)
}

type tModalMethod uint8

const (
	modalMethodOnAttachedToNavigator tModalMethod = iota
	modalMethodOnResume
	modalMethodHandle
	modalMethodLayout
	modalMethodOnDismiss
)

func combineModalMethods(methods []tModalMethod) string {
	methodNames := make([]string, len(methods))
	for i := range methods {
		methodNames[i] = methods[i].String()
	}
	return strings.Join(methodNames, ",")
}

func (tpm tModalMethod) String() string {
	switch tpm {
	case modalMethodOnAttachedToNavigator:
		return "OnAttachedToNavigator()"
	case modalMethodOnResume:
		return "OnResume()"
	case modalMethodHandle:
		return "Handle()"
	case modalMethodLayout:
		return "Layout()"
	case modalMethodOnDismiss:
		return "OnDismiss()"
	default:
		return fmt.Sprintf("UnknownMethod %d", tpm)
	}
}

type testModal struct {
	id               string
	log              func(format string, args ...interface{})
	mtx              sync.Mutex
	parentNav        PageNavigator
	methodsCallStack []tModalMethod
}

func newTestModal(id string, logFn func(format string, args ...interface{})) *testModal {
	return &testModal{
		id:  id,
		log: logFn,
	}
}

func (tModal *testModal) calledMethods() []tModalMethod {
	tModal.mtx.Lock()
	defer tModal.mtx.Unlock()
	methods := tModal.methodsCallStack
	tModal.methodsCallStack = nil
	tModal.log("cleared called methods on %s", tModal.id)
	return methods
}

func (tModal *testModal) calledMethod(tpm tModalMethod) {
	tModal.mtx.Lock()
	defer tModal.mtx.Unlock()
	tModal.methodsCallStack = append(tModal.methodsCallStack, tpm)
	tModal.log("recorded %s.%s", tModal.id, tpm)
}

func (tModal *testModal) ID() string {
	// No need to track calls to the ID() method, doesn't do anything special.
	return tModal.id
}

func (tModal *testModal) OnAttachedToNavigator(navigator PageNavigator) {
	tModal.calledMethod(modalMethodOnAttachedToNavigator)
	tModal.parentNav = navigator
}

func (tModal *testModal) OnResume() {
	tModal.calledMethod(modalMethodOnResume)
}

func (tModal *testModal) Handle() {
	tModal.calledMethod(modalMethodHandle)
}

func (tModal *testModal) Layout(layout.Context) layout.Dimensions {
	tModal.calledMethod(modalMethodLayout)
	return layout.Dimensions{}
}

func (tModal *testModal) OnDismiss() {
	tModal.calledMethod(modalMethodOnDismiss)
}

// testGiouiWindow is a simplified version of gioui.org/app.Window.
type testGiouiWindow struct {
	topModalGetter    func() Modal
	currentPageGetter func() Page
}

// invalidate is a simplified version of gioui.org/app.Window.Invalidate(). The
// real version sends a gioui.org/io/system.FrameEvent which consumers use to
// render the window display onscreen. This simplified version just calls the
// appropriate methods on the page to be displayed.
func (tgw *testGiouiWindow) invalidate() {
	if topModal := tgw.topModalGetter(); topModal != nil {
		topModal.Handle()
		topModal.Layout(layout.Context{})
	} else if currentPage := tgw.currentPageGetter(); currentPage != nil {
		currentPage.HandleUserInteractions()
		currentPage.Layout(layout.Context{})
	}
}

func TestSimpleWindowNavigator(mainT *testing.T) {
	tWindow := &testGiouiWindow{}
	windowNavigator := NewSimpleWindowNavigator(tWindow.invalidate)
	tWindow.topModalGetter = windowNavigator.TopModal
	tWindow.currentPageGetter = windowNavigator.CurrentPage

	tests := []struct {
		methodToTest string // display, closeCurrentPage
		id           string // = newPageID (display tests), finalPageID (closeCurrentPage tests), modalID (modal tests).

		expectEmptyStackBefore, expectEmptyStackAfter bool
	}{
		{methodToTest: "closeCurrentPage", expectEmptyStackBefore: true, expectEmptyStackAfter: true},
		{methodToTest: "display", id: generateTestID("page", 1), expectEmptyStackBefore: true},
		{methodToTest: "display", id: generateTestID("page", 1)},                               // test displaying duplicate page(1)
		{methodToTest: "closeCurrentPage", expectEmptyStackAfter: true},                        // page(1) should close, leaving nothing
		{methodToTest: "display", id: generateTestID("page", 1), expectEmptyStackBefore: true}, // redisplay page(1)
		{methodToTest: "display", id: generateTestID("page", 2)},
		{methodToTest: "display", id: generateTestID("page", 2)},          // test displaying duplicate page(2)
		{methodToTest: "closeCurrentPage", id: generateTestID("page", 1)}, // page(2) should close, leaving page(1)
		{methodToTest: "display", id: generateTestID("page", 2)},          // re-displaying page(2)
		{methodToTest: "display", id: generateTestID("page", 3)},
		{methodToTest: "closeCurrentPage", id: generateTestID("page", 2)}, // page(3) should close, leaving page(2)
		{methodToTest: "closePagesAfter", id: generateTestID("page", 1)},  // page(2,3) should close, leaving page(1)
		{methodToTest: "closePagesAfter"},                                 // no id specified to close after, no page should close
		{methodToTest: "closeAllPages", expectEmptyStackAfter: true},
		{methodToTest: "display", id: generateTestID("page", 1), expectEmptyStackBefore: true},
		{methodToTest: "display", id: generateTestID("page", 2)},
		{methodToTest: "display", id: generateTestID("page", 3)},
		{methodToTest: "display", id: generateTestID("page", 1)}, // re-displaying page(1) should kill previous page(1)
		{methodToTest: "clearStackAndDisplay", id: generateTestID("page", 1)},
		{methodToTest: "showModal", id: generateTestID("modal", 1), expectEmptyStackBefore: true},
		{methodToTest: "showModal", id: generateTestID("modal", 1)}, // duplicates allowed
		{methodToTest: "showModal", id: generateTestID("modal", 2)},
		{methodToTest: "showModal", id: generateTestID("modal", 1)},
		{methodToTest: "showModal", id: generateTestID("modal", 2)},
		{methodToTest: "dismissModal", id: generateTestID("modal", 3)}, // does not exist
		{methodToTest: "dismissModal", id: generateTestID("modal", 1)},
		{methodToTest: "dismissModal", id: generateTestID("modal", 1)},
		{methodToTest: "dismissModal", id: generateTestID("modal", 2)},
		{methodToTest: "dismissModal", id: generateTestID("modal", 2)},
		{methodToTest: "dismissModal", id: generateTestID("modal", 1), expectEmptyStackAfter: true}, // last modal dismissed
	}

	for _, tt := range tests {
		var testName string
		var runTest func(*testing.T)

		switch tt.methodToTest {
		case "display":
			testName = fmt.Sprintf("display %s", tt.id)
			runTest = func(t *testing.T) {
				testWindowNavigatorDisplay(t, tt.id, windowNavigator)
			}

		case "closeCurrentPage":
			testName = fmt.Sprintf("closeCurrentPage %s", windowNavigator.CurrentPageID())
			runTest = func(t *testing.T) {
				var finalPageID *string
				if !tt.expectEmptyStackAfter {
					finalPageID = &tt.id
				}
				testWindowNavigatorCloseCurrentPage(t, finalPageID, windowNavigator)
			}

		case "closePagesAfter":
			testName = fmt.Sprintf("closePagesAfter %s", tt.id)
			if tt.id == "" {
				testName = "closePagesAfter non-existing page"
			}
			runTest = func(t *testing.T) {
				var pagesToBeClosed []Page
				if tt.id != "" && tt.id != windowNavigator.CurrentPageID() {
					pagesToBeClosed = windowNavigator.subPages.pagesAfter(&tt.id)
				}
				testWindowNavigatorClosePagesAfter(t, tt.id, pagesToBeClosed, windowNavigator)
			}

		case "closeAllPages":
			testName = "closeAllPages"
			tt.expectEmptyStackAfter = true // ensures that all pages are closed
			runTest = func(t *testing.T) {
				pagesToBeClosed := windowNavigator.subPages.pagesAfter(nil) // all pages
				testWindowNavigatorCloseAllPages(t, pagesToBeClosed, windowNavigator)
			}

		case "clearStackAndDisplay":
			testName = fmt.Sprintf("clearStackAndDisplay %s", tt.id)
			runTest = func(t *testing.T) {
				pagesToBeClosed := windowNavigator.subPages.pagesAfter(nil) // all pages
				testWindowNavigatorClearStackAndDisplay(t, tt.id, pagesToBeClosed, windowNavigator)
			}

		case "showModal":
			testName = fmt.Sprintf("showModal %s", tt.id)
			runTest = func(t *testing.T) {
				validateModalCountDifference(t, windowNavigator, 1, func() {
					testWindowNavigatorShowModal(t, tt.id, windowNavigator)
				})
			}

		case "dismissModal":
			testName = fmt.Sprintf("dismissModal %s", tt.id)
			runTest = func(t *testing.T) {
				var modalToDismiss Modal
				var modalCountDiff int
				windowNavigator.modalMutex.Lock()
				for i := len(windowNavigator.modals) - 1; i >= 0; i-- {
					modal := windowNavigator.modals[i]
					if modal.ID() == tt.id {
						modalToDismiss = modal
						modalCountDiff = -1
						break
					}
				}
				windowNavigator.modalMutex.Unlock()

				validateModalCountDifference(t, windowNavigator, modalCountDiff, func() {
					testWindowNavigatorDismissModal(t, tt.id, modalToDismiss, windowNavigator)
				})
			}

		default:
			mainT.Fatalf("unexpected methodToTest: %q", tt.methodToTest)
		}

		checkNavigatorStack := func(t *testing.T, isModalNav, expectEmpty bool, beforeOrAfter string) {
			var emptyStack bool
			if isModalNav {
				emptyStack = windowNavigator.TopModal() == nil
			} else {
				emptyStack = windowNavigator.CurrentPage() == nil
			}
			if emptyStack != expectEmpty {
				pageOrModal := "page"
				if isModalNav {
					pageOrModal = "modal"
				}
				t.Fatalf("expected empty %s stack %s navigation: %t, found empty stack %t",
					pageOrModal, beforeOrAfter, expectEmpty, emptyStack)
			}
		}

		mainT.Run(testName, func(t *testing.T) {
			// Update log function for pages and modals in the navigator's stack.
			allPages := windowNavigator.subPages.pagesAfter(nil)
			for _, page := range allPages {
				if page, ok := page.(*testPage); ok {
					page.log = t.Logf
				}
			}
			windowNavigator.modalMutex.Lock()
			for _, modal := range windowNavigator.modals {
				if modal, ok := modal.(*testModal); ok {
					modal.log = t.Logf
				}
			}
			windowNavigator.modalMutex.Unlock()

			isModalNav := strings.Contains(tt.id, "modal")
			checkNavigatorStack(t, isModalNav, tt.expectEmptyStackBefore, "before")
			runTest(t)
			checkNavigatorStack(t, isModalNav, tt.expectEmptyStackAfter, "after")
		})
		println() // just for prettier display
	}
}

func testWindowNavigatorDisplay(t *testing.T, newPageID string, windowNavigator WindowNavigator) {
	pageBeforeDisplay := windowNavigator.CurrentPage()
	newPage := newTestPage(newPageID, t.Logf)
	windowNavigator.Display(newPage)

	if pageBeforeDisplay != nil && pageBeforeDisplay.ID() == newPage.ID() {
		// new page should not be displayed
		if newPage.parentNav != nil {
			t.Fatalf("found wrong parentNav for new page that wasn't displayed: %T", newPage.parentNav)
		}

		// Confirm that no new methods are called on the previously displayed page.
		previousPage, ok := pageBeforeDisplay.(*testPage)
		if !ok {
			t.Fatalf("previous page of unexpected type %T", previousPage)
		}
		previousPageCalledMethods := previousPage.calledMethods()
		if len(previousPageCalledMethods) != 0 {
			t.Fatalf("unexpected methods called on the previous page: %s",
				combinePageMethods(previousPageCalledMethods))
		}
		return
	}

	validateDisplayedPage(t, newPage, false)
	if pageBeforeDisplay != nil {
		validateDismissedPage(t, pageBeforeDisplay, false)
	}
}

func testWindowNavigatorCloseCurrentPage(t *testing.T, finalPageID *string, windowNavigator WindowNavigator) {
	pageToClose := windowNavigator.CurrentPage()
	windowNavigator.CloseCurrentPage()
	finalPage := windowNavigator.CurrentPage()

	// If there was no page to close, there should be no final page.
	if pageToClose == nil {
		if finalPageID != nil {
			t.Fatalf("bad test: no page to close but expected final page %q", *finalPageID)
		}
		if finalPage != nil {
			t.Fatalf("unexpected final page: %q", finalPage.ID())
		}
		return // nothing further to check, both pageToClose and finalPage are nil.
	}

	// Check that the closed page was properly dimissed.
	validateDismissedPage(t, pageToClose, true)

	// If no final page was expected, confirm that no final page was found.
	if finalPageID == nil {
		if finalPage != nil {
			t.Fatalf("unexpected final page: %q", finalPage.ID())
		}
		return // nothing further to check, no final page expected, no final page found.
	}

	// Confirm that final page is as expected.
	if finalPage == nil {
		t.Fatalf("expected final page %q but found none", *finalPageID)
	}
	if *finalPageID != finalPage.ID() {
		t.Fatalf("expected final page %q but found %q", *finalPageID, finalPage.ID())
	}
	if finalPage, ok := finalPage.(*testPage); ok {
		validateDisplayedPage(t, finalPage, true)
	} else {
		t.Fatalf("found final page of unexpected type %T", finalPage)
	}
}

func testWindowNavigatorClosePagesAfter(t *testing.T, finalPageID string, pagesToBeClosed []Page, windowNavigator WindowNavigator) {
	pageBeforeClosing := windowNavigator.CurrentPage()
	windowNavigator.ClosePagesAfter(finalPageID)
	finalPage := windowNavigator.CurrentPage()

	// If there was no page to close, there should be no final page.
	if pageBeforeClosing == nil {
		if len(pagesToBeClosed) > 0 {
			t.Fatalf("bad test: no page to close but expected %d closed pages", len(pagesToBeClosed))
		}
		if finalPageID != "" {
			t.Fatalf("bad test: no page to close but expected final page %q", finalPageID)
		}
		if finalPage != nil {
			t.Fatalf("unexpected final page: %q", finalPage.ID())
		}
		return // nothing further to check, both pageBeforeClosing and finalPage are nil.
	}

	// Closing pages after an invalid page ID or same ID as the current page
	// should do nothing. Final page ID should be unchanged.
	if finalPageID == "" || pageBeforeClosing.ID() == finalPageID {
		if len(pagesToBeClosed) > 0 {
			t.Fatalf("bad test: no page should be closed but expected %d closed pages", len(pagesToBeClosed))
		}
		if finalPage.ID() != pageBeforeClosing.ID() {
			t.Fatalf("expected final page ID to be %s, got %s", pageBeforeClosing.ID(), finalPage.ID())
		}
		return
	}

	// Confirm that there is a final page and it is as expected.
	if finalPage == nil {
		t.Fatalf("expected final page %q but found none", finalPageID)
	}
	// The final page cannot be the same page as the previous top page.
	if finalPage.ID() == pageBeforeClosing.ID() {
		t.Fatalf("previous top page %s not closed", pageBeforeClosing.ID())
	}
	// Confirm that final page ID is as expected.
	if finalPageID != finalPage.ID() {
		t.Fatalf("expected final page %q but found %q", finalPageID, finalPage.ID())
	}
	// Validate the final page.
	if finalPage, ok := finalPage.(*testPage); ok {
		validateDisplayedPage(t, finalPage, true)
	} else {
		t.Fatalf("found final page of unexpected type %T", finalPage)
	}

	// Confirm that the closed page(s) were properly dimissed.
	for _, page := range pagesToBeClosed {
		validateDismissedPage(t, page, true)
	}
}

func testWindowNavigatorCloseAllPages(t *testing.T, pagesToBeClosed []Page, windowNavigator WindowNavigator) {
	windowNavigator.CloseAllPages()
	for _, page := range pagesToBeClosed {
		validateDismissedPage(t, page, true)
	}
}

func testWindowNavigatorClearStackAndDisplay(t *testing.T, newPageID string, pagesToBeClosed []Page, windowNavigator WindowNavigator) {
	newPage := newTestPage(newPageID, t.Logf)
	windowNavigator.ClearStackAndDisplay(newPage)

	validateDisplayedPage(t, newPage, false)
	for _, page := range pagesToBeClosed {
		validateDismissedPage(t, page, true)
	}
}

func validateModalCountDifference(t *testing.T, windowNavigator *SimpleWindowNavigator, expectedDiff int, modifyModals func()) {
	countModals := func() int {
		windowNavigator.modalMutex.Lock()
		defer windowNavigator.modalMutex.Unlock()
		return len(windowNavigator.modals)
	}
	initialCount := countModals()
	modifyModals()
	finalCount := countModals()
	if finalCount-initialCount != expectedDiff {
		t.Fatalf("found %d modals after test (previously %d) instead of %d",
			finalCount, initialCount, initialCount+expectedDiff)
	}
}

func testWindowNavigatorShowModal(t *testing.T, modalID string, windowNavigator WindowNavigator) {
	modalBeforeNew := windowNavigator.TopModal()
	newModal := newTestModal(modalID, t.Logf)
	windowNavigator.ShowModal(newModal)

	// Confirm that the correct parentNav is set.
	switch newModal.parentNav.(type) {
	case WindowNavigator:
	default:
		t.Fatalf("found wrong parentNav for displayed modal: %T", newModal.parentNav)
	}

	// Confirm that the right methods are called on the new modal.
	expectedMethods := []tModalMethod{
		modalMethodOnAttachedToNavigator,
		modalMethodOnResume,
		modalMethodHandle,
		modalMethodLayout,
	}
	newModalCalledMethods := newModal.calledMethods()
	if len(newModalCalledMethods) != len(expectedMethods) {
		t.Fatalf("%s instead of %s methods were called on the new modal",
			combineModalMethods(newModalCalledMethods), combineModalMethods(expectedMethods))
	}
	for i, expectedMethod := range expectedMethods {
		foundMethod := newModalCalledMethods[i]
		if foundMethod != expectedMethod {
			t.Fatalf("expected %s method called on new modal to be %s, but found %s",
				numerate(i+1), expectedMethod, foundMethod)
		}
	}

	// Confirm that no new methods are called on the previously displayed modal.
	if modalBeforeNew != nil {
		previousModal, ok := modalBeforeNew.(*testModal)
		if !ok {
			t.Fatalf("previous modal of unexpected type %T", previousModal)
		}
		previousModalCalledMethods := previousModal.calledMethods()
		if len(previousModalCalledMethods) != 0 {
			t.Fatalf("unexpected methods called on the previous modal: %s",
				combineModalMethods(previousModalCalledMethods))
		}
	}
}

func testWindowNavigatorDismissModal(t *testing.T, modalID string, modalToBeClosed Modal, windowNavigator WindowNavigator) {
	windowNavigator.DismissModal(modalID)

	// Confirm that the correct parentNav is set for the new top modal
	// and that the right methods are called on it.
	if newTopModal := windowNavigator.TopModal(); newTopModal != nil {
		newTopModal, ok := newTopModal.(*testModal)
		if !ok {
			t.Fatalf("top modal of unexpected type %T", newTopModal)
		}
		switch newTopModal.parentNav.(type) {
		case WindowNavigator:
		default:
			t.Fatalf("found wrong parentNav for displayed modal: %T", newTopModal.parentNav)
		}

		expectedMethods := []tModalMethod{
			modalMethodHandle,
			modalMethodLayout,
		}
		if modalToBeClosed == nil {
			expectedMethods = nil
		}
		newModalCalledMethods := newTopModal.calledMethods()
		if len(newModalCalledMethods) != len(expectedMethods) {
			t.Fatalf("%s instead of %s methods were called on the new modal",
				combineModalMethods(newModalCalledMethods), combineModalMethods(expectedMethods))
		}
		for i, expectedMethod := range expectedMethods {
			foundMethod := newModalCalledMethods[i]
			if foundMethod != expectedMethod {
				t.Fatalf("expected %s method called on new modal to be %s, but found %s",
					numerate(i+1), expectedMethod, foundMethod)
			}
		}
	}

	if modalToBeClosed == nil {
		return
	}

	// Confirm that the right methods were called on the dismissed modal.
	dismissedModal, ok := modalToBeClosed.(*testModal)
	if !ok {
		t.Fatalf("dismissed modal of unexpected type %T", modalToBeClosed)
	}
	expectedMethods := []tModalMethod{modalMethodOnDismiss}
	newPageCalledMethods := dismissedModal.calledMethods()
	if len(newPageCalledMethods) != len(expectedMethods) {
		t.Fatalf("%s instead of %s methods were called on the dismissed modal",
			combineModalMethods(newPageCalledMethods), combineModalMethods(expectedMethods))
	}
	for i, expectedMethod := range expectedMethods {
		foundMethod := newPageCalledMethods[i]
		if foundMethod != expectedMethod {
			t.Fatalf("expected %s method called on dismissed modal to be %s, but found %s",
				numerate(i+1), expectedMethod, foundMethod)
		}
	}
}

func validateDisplayedPage(t *testing.T, newPage *testPage, isRedisplay bool) {
	// Confirm that the correct parentNav is set.
	switch newPage.parentNav.(type) {
	case WindowNavigator:
	default:
		t.Fatalf("found wrong parentNav for displayed page: %T", newPage.parentNav)
	}

	// Confirm that the right methods are called on the new page.
	expectedMethods := []tPageMethod{
		pageMethodOnAttachedToNavigator,
		pageMethodOnNavigatedTo,
		pageMethodHandleUserInteractions,
		pageMethodLayout,
	}
	if isRedisplay {
		expectedMethods = expectedMethods[1:] // OnAttachedToNavigator only called on first display
	}
	newPageCalledMethods := newPage.calledMethods()
	if len(newPageCalledMethods) != len(expectedMethods) {
		t.Fatalf("%s instead of %s methods were called on the new page",
			combinePageMethods(newPageCalledMethods), combinePageMethods(expectedMethods))
	}
	for i, expectedMethod := range expectedMethods {
		foundMethod := newPageCalledMethods[i]
		if foundMethod != expectedMethod {
			t.Fatalf("expected %s method called on new page to be %s, but found %s",
				numerate(i+1), expectedMethod, foundMethod)
		}
	}
}

func validateDismissedPage(t *testing.T, page Page, wasClosed bool) {
	dismissedPage, ok := page.(*testPage)
	if !ok {
		t.Fatalf("dismissed page of unexpected type %T", dismissedPage)
	}
	expectedMethods := []tPageMethod{pageMethodOnNavigatedFrom}
	if wasClosed {
		expectedMethods = append(expectedMethods, pageMethodOnClosed)
	}
	closedPageCalledMethods := dismissedPage.calledMethods()
	if len(closedPageCalledMethods) != len(expectedMethods) {
		t.Fatalf("%s instead of %s methods were called on closed page %s",
			combinePageMethods(closedPageCalledMethods), combinePageMethods(expectedMethods), dismissedPage.ID())
	}
	for i, expectedMethod := range expectedMethods {
		if closedPageCalledMethods[i] != expectedMethod {
			t.Fatalf("expected %s method to be called on closed page, but found %s",
				expectedMethod, closedPageCalledMethods[0])
		}
	}
}

func generateTestID(pageOrModal string, i int) string {
	return fmt.Sprintf("test%s-%d", pageOrModal, i)
}

func numerate(n int) string {
	nStr := strconv.Itoa(n)
	lastDigit := nStr[len(nStr)-1]
	switch lastDigit {
	case '1':
		return nStr + "st"
	case '2':
		return nStr + "nd"
	case '3':
		return nStr + "rd"
	default:
		return nStr + "th"
	}
}
