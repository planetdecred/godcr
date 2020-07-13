// SPDX-License-Identifier: Unlicense OR MIT

package decredmaterial

import (
	_ "image/png" //makes png images a decodable format

	"gioui.org/widget"
)

type Icon struct {
	*widget.Icon
}

// NewIcon returns a new Icon from IconVG data.
func NewIcon(data []byte) (*Icon, error) {
	icon, err := widget.NewIcon(data)
	if err != nil {
		return nil, err
	}
	return &Icon{icon}, nil
}
