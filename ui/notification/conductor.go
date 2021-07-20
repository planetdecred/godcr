package notification

import (
	"sync"
)

type Conductor struct {
	queue     chan Notification
	active    Notification
	activeMtx sync.Mutex
}

func NewConductor() *Conductor {
	c := &Conductor{
		queue: make(chan Notification, 5),
	}
	go c.conduct()

	return c
}

func (c *Conductor) conduct() {
	for ntf := range c.queue {
		c.activeMtx.Lock()
		c.active = ntf
		c.activeMtx.Unlock()
	}
}

func (c *Conductor) GetActiveNotification() Notification {
	c.activeMtx.Lock()
	defer c.activeMtx.Unlock()

	return c.active
}

func (c *Conductor) NewNotifier() Notifier {
	return Notifier{
		queue: c.queue,
	}
}
