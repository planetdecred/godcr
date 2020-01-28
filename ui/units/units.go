// Package units provides unit values used across the app
package units

import (
	"gioui.org/unit"
)

var (
	// Label is the unit for the app labels
	Label = unit.Dp(50)

	// FlexInset is the unit for flex insets
	FlexInset = unit.Dp(50)

	// Padding is the unit for uniform padding for widgets in a column
	Padding = unit.Dp(5)

	// ContainerPadding is the unit for uniform padding for content in a container
	ContainerPadding = unit.Dp(20)

	// PageMarginTop is the unit for top margin of pages
	PageMarginTop = unit.Dp(50)

	// ColumnMargin is the unit for top margins of columns of the overview page
	ColumnMargin = unit.Dp(30)

	// TransactionBalanceMain is the unit for the main balance text size on a transaction row
	TransactionBalanceMain = unit.Dp(14)

	// TransactionBalanceMain is the unit for the main balance text size on a transaction row
	TransactionBalanceSub = unit.Dp(10)

	// TransactionsRowMargin is the unit for uniform spacing between widgets in a recent transaction row
	TransactionsRowMargin = unit.Dp(10)

	// SyncBoxPadding is the unit for uniform padding of sync wallet boxes
	SyncBoxPadding = unit.Dp(10)

	// NoPadding is the unit for applying zero padding for widgets in a list layout
	NoPadding = unit.Dp(0)

	// Padding1 is a unit value of one
	Padding1 = unit.Dp(1)

	// WalletSyncBoxWidthMin is the unit for the minimum width of a sync wallet box
	WalletSyncBoxWidthMin = unit.Dp(300)

	// WalletSyncBoxHeightMin is the unit for the minimum height of a sync wallet box
	WalletSyncBoxHeightMin = unit.Dp(90)

	// WalletSyncBoxContentWidth is the unit for the maximum and minimum width of the
	// contents of a wallet sync box
	WalletSyncBoxContentWidth = unit.Dp(280)

	// SyncButtonTextSize is the unit for the text size of sync button
	SyncButtonTextSize = unit.Dp(10)

	// DefaultButtonRadius is the unit for corner radius of material buttons
	DefaultButtonRadius = unit.Dp(4)
)
