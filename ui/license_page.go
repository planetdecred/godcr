package ui

import (
	"gioui.org/layout"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
)

const PageLicense = "License"

const License = `ISC License

Copyright (c) 2018-2021, Raedah Group

Permission to use, copy, modify, and/or distribute this software for any
purpose with or without fee is hereby granted, provided that the above
copyright notice and this permission notice appear in all copies.

THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.`

type licensePage struct {
	common        *pageCommon
	pageContainer layout.List

	backButton decredmaterial.IconButton
}

func LicensePage(common *pageCommon) Page {
	pg := &licensePage{
		common:        common,
		pageContainer: layout.List{Axis: layout.Vertical},
	}
	pg.backButton, _ = common.SubPageHeaderButtons()

	return pg
}

func (pg *licensePage) OnResume() {

}

//main page layout
func (pg *licensePage) Layout(gtx layout.Context) layout.Dimensions {
	common := pg.common
	d := func(gtx C) D {
		page := SubPage{
			title:      "License",
			backButton: pg.backButton,
			back: func() {
				pg.common.changePage(PageAbout)
			},
			body: func(gtx C) D {
				return common.theme.Card().Layout(gtx, func(gtx C) D {
					return layout.UniformInset(values.MarginPadding25).Layout(gtx, func(gtx C) D {
						licenseText := common.theme.Body1(License)
						licenseText.Color = common.theme.Color.Gray
						return layout.Inset{Bottom: values.MarginPadding20}.Layout(gtx, licenseText.Layout)
					})
				})
			},
		}
		return common.SubPageLayout(gtx, page)
	}
	return common.UniformPadding(gtx, d)
}

func (pg *licensePage) handle() {

}

func (pg *licensePage) onClose() {}
