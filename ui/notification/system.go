package notification

import (
	"fmt"
	"os"
	"path"
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

func getAbsolutePath() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("error getting executable path: %s", err.Error())
	}

	exSym, err := filepath.EvalSymlinks(ex)
	if err != nil {
		return "", fmt.Errorf("error getting filepath after evaluating sym links")
	}

	return path.Dir(exSym), nil
}
