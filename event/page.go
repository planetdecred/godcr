package event

// Nav represents an event for the window to change it's current page
type Nav struct {
	Current string
	Next    string
}
