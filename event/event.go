// Package event provides an types for passing events between the app's
// components.
package event

// Event is a convenient type for representing events.
// An error can also be passed as an Event.
type Event interface{}

// Duplex is a stucture for constraning event communication to two
// directionally-constrained channels so as to reduce the possibility
// of a deadlock
type Duplex struct {
	Send    (chan<- Event)
	Receive (<-chan Event)
}
