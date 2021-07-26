package notification

import (
	"path/filepath"

	"github.com/gen2brain/beeep"
)

const (
	icon  = "ui/assets/decredicons/qrcodeSymbol.png"
	title = "Godcr"
)

type SystemNotification struct {
	iconPath string
	message  string
}

func NewSystemNotification() *SystemNotification {
	absolutePath, err := getAbsolutePath()
	if err != nil {
		log.Error(err.Error())
	}

	return &SystemNotification{
		iconPath: filepath.Join(absolutePath, icon),
	}
}

func (s *SystemNotification) Notify(message string) {
	err := beeep.Notify(title, message, s.iconPath)
	if err != nil {
		log.Info("could not initiate desktop notification, reason:", err.Error())
	}
}
