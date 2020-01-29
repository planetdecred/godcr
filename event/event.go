// Package event provides an types for passing events between the app's
// components.
package event

import (
	"errors"
)

// Event is a convenient type for representing events.
// An error can also be passed as an Event.
type Event interface{}

// Duplex is a structure for constraning event communication to two
// directionally-constrained channels so as to reduce the possibility
// of a deadlock
type Duplex struct {
	Send    (chan<- Event)
	Receive (<-chan Event)
}

// DuplexBase is a struct for creating a Duplex and it's reverse
type DuplexBase struct {
	A, B chan Event
}

// NewDuplexBase creates a new DuplexBase
func NewDuplexBase() DuplexBase {
	return DuplexBase{
		A: make(chan Event, 2),
		B: make(chan Event, 2),
	}
}

// Duplex returns the Duplex
func (dup DuplexBase) Duplex() Duplex {
	return Duplex{
		Send:    dup.B,
		Receive: dup.A,
	}
}

// Reverse returns the reversed Duplex
func (dup DuplexBase) Reverse() Duplex {
	return Duplex{
		Send:    dup.A,
		Receive: dup.B,
	}
}

var (
	// ErrQueueUnderflow is returned when ArgumentQueue is empty and a Pop is requested
	ErrQueueUnderflow = errors.New("no more arguments")
	// ErrInvalidPop is returned when the current ArgumentQueue item is not of the type requested
	ErrInvalidPop = errors.New("current argument cannot be asserted to requested type")
)

// ArgumentQueue is a structure for handling data passed through events
type ArgumentQueue struct {
	Queue []interface{}
}

func (queue *ArgumentQueue) pop() (interface{}, error) {
	q := queue.Queue
	if len(q) == 0 {
		return nil, ErrQueueUnderflow
	}
	p := q[0]
	queue.Queue = q[1:]
	return p, nil

}

// PopString pops a string from the queue.
// It returns an error when the queue is empty and
// when the current item is not a string
func (queue *ArgumentQueue) PopString() (string, error) {
	s, err := queue.pop()
	if err != nil {
		return "", err
	}
	str, ok := s.(string)
	if !ok {
		return "", ErrInvalidPop
	}
	return str, nil
}

// PopInt pops a string from the queue.
// It returns an error when the queue is empty and
// when the current item is not an int
func (queue *ArgumentQueue) PopInt() (int, error) {
	i, err := queue.pop()
	if err != nil {
		return 0, err
	}
	in, ok := i.(int)
	if !ok {
		return 0, ErrInvalidPop
	}
	return in, nil
}
