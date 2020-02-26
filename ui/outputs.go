package ui

import "gioui.org/widget/material"

type outputs struct {
	spendingPassword, matchSpending material.Editor
	toLanding, toWallets            material.IconButton
	icons                           struct {
		add, check *material.Icon
	}
}
