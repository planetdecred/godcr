package notification

import "time"

type Notifier struct {
	queue chan<- Notification
}

func (n Notifier) Notify(text string, t NotificationType) {
	n.queue <- Notification{
		NotificationType: t,
		text:             text,
		created:          time.Now(),
	}
}
