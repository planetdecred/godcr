package notification

import (
	"time"
)

type NotificationType int

const (
	Success NotificationType = iota
	Warning
	Error
)

type Notification struct {
	NotificationType
	text    string
	created time.Time
}
