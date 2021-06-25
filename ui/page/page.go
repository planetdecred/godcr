package page

import (
	"gioui.org/gesture"
	"gioui.org/unit"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/modal"
	"image/color"

	"gioui.org/layout"
	"gioui.org/widget"

	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
)

type navHandler struct {
	clickable     *widget.Clickable
	image         *widget.Image
	imageInactive *widget.Image
	page          string
}

type walletAccount struct {
	evt          *gesture.Click
	walletIndex  int
	accountIndex int
	accountName  string
	totalBalance string
	spendable    string
	number       int32
}

type wallectAccountOption struct {
	selectSendAccount           map[int][]walletAccount
	selectReceiveAccount        map[int][]walletAccount
	selectPurchaseTicketAccount map[int][]walletAccount
}


type WalletAccountSelector struct {
	title                     string
	walletAccount             decredmaterial.Modal
	walletsList, accountsList layout.List
	isWalletAccountModalOpen  bool
	isWalletAccountInfo       bool
	walletAccounts            *wallectAccountOption
	sendAccountBtn            *widget.Clickable
	receivingAccountBtn       *widget.Clickable
	purchaseTicketAccountBtn  *widget.Clickable
	sendOption                string
	walletInfoButton          decredmaterial.IconButton

	selectedSendAccount,
	selectedSendWallet,
	selectedReceiveAccount,
	selectedReceiveWallet,
	selectedPurchaseTicketAccount,
	selectedPurchaseTicketWallet int
}

type (
	C = layout.Context
	D = layout.Dimensions
)

// Container is simply a wrapper for the Inset type. Its purpose is to differentiate the use of an inset as a padding or
// margin, making it easier to visualize the structure of a layout when reading UI code.
type Container struct {
	padding layout.Inset
}

func (c Container) Layout(gtx layout.Context, w layout.Widget) layout.Dimensions {
	return c.padding.Layout(gtx, w)
}

var (
	MaxWidth = unit.Dp(800)
)

type SubPage struct {
	*load.Load
	title        string
	subTitle     string
	walletName   string
	back         func()
	body         layout.Widget
	infoTemplate string
	extraItem    *widget.Clickable
	extra        layout.Widget
	extraText    string
	handleExtra  func()
}

func (sp *SubPage) Layout(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Bottom: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
				return sp.Header(gtx)
			})
		}),
		layout.Rigid(sp.body),
	)
}

func (sp *SubPage) Header(gtx layout.Context) layout.Dimensions {
	sp.EventHandler()

	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Right: values.MarginPadding20}.Layout(gtx, sp.SubPageBackButton.Layout)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if sp.subTitle == "" {
				return sp.Theme.H6(sp.title).Layout(gtx)
			}
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(sp.Theme.H6(sp.title).Layout),
				layout.Rigid(sp.Theme.Body1(sp.subTitle).Layout),
			)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if sp.walletName != "" {
				return layout.Inset{Left: values.MarginPadding5, Top: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
					return decredmaterial.Card{
						Color: sp.Theme.Color.Surface,
					}.Layout(gtx, func(gtx C) D {
						return layout.UniformInset(values.MarginPadding2).Layout(gtx, func(gtx C) D {
							walletText := sp.Theme.Caption(sp.walletName)
							walletText.Color = sp.Theme.Color.Gray
							return walletText.Layout(gtx)
						})
					})
				})
			}
			return layout.Dimensions{}
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return layout.E.Layout(gtx, func(gtx C) D {
				if sp.infoTemplate != "" {
					return sp.SubPageInfoButton.Layout(gtx)
				} else if sp.extraItem != nil {
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							if sp.extraText != "" {
								return layout.Inset{Right: values.MarginPadding10, Top: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
									text := sp.Theme.Caption(sp.extraText)
									text.Color = sp.Theme.Color.DeepBlue
									return text.Layout(gtx)
								})
							}
							return layout.Dimensions{}
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return decredmaterial.Clickable(gtx, sp.extraItem, sp.extra)
						}),
					)
				}
				return layout.Dimensions{}
			})
		}),
	)
}

func (sp *SubPage) SplitLayout(gtx layout.Context) layout.Dimensions {
	card := sp.Theme.Card()
	card.Color = color.NRGBA{}
	return card.Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D { return sp.Header(gtx)}),
			layout.Rigid(sp.body),
		)
	})
}

func (sp *SubPage) EventHandler() {
	if sp.SubPageInfoButton.Button.Clicked() {
		modal.NewInfoModal(sp.Load).
			Title(sp.title).
			SetupWithTemplate(sp.infoTemplate).
			NegativeButton("Got it", func() {}).Show()
	}

	if sp.SubPageBackButton.Button.Clicked() {
		sp.back()
	}

	if sp.extraItem != nil && sp.extraItem.Clicked() {
		sp.handleExtra()
	}
}

func uniformPadding(gtx layout.Context, body layout.Widget) layout.Dimensions {
	width := gtx.Constraints.Max.X

	padding := values.MarginPadding24

	if (width - 2*gtx.Px(padding)) > gtx.Px(MaxWidth) {
		paddingValue := float32(width-gtx.Px(MaxWidth)) / 2
		padding = unit.Px(paddingValue)
	}

	return layout.Inset{
		Top:    values.MarginPadding24,
		Right:  padding,
		Bottom: values.MarginPadding24,
		Left:   padding,
	}.Layout(gtx, body)
}

func loadPages(l *load.Load) map[string]load.Page {

	iconColor := l.Theme.Color.Gray3

	zeroInset := layout.UniformInset(values.MarginPadding0)
	l.SubPageBackButton.Color, l.SubPageInfoButton.Color = iconColor, iconColor

	m25 := values.MarginPadding25
	l.SubPageBackButton.Size, l.SubPageInfoButton.Size = m25, m25
	l.SubPageBackButton.Inset, l.SubPageInfoButton.Inset = zeroInset, zeroInset

	pages := make(map[string]load.Page)

	pages[Wallet] = WalletPage(l)
	pages[More] = MorePage(l)
	pages[CreateRestore] = CreateRestorePage(l)
	pages[Receive] = ReceivePage(l)
	pages[Send] = SendPage(l)
	pages[SignMessage] = SignMessagePage(l)
	pages[VerifyMessage] = VerifyMessagePage(l)
	pages[SeedBackup] = BackupPage(l)
	pages[Settings] = SettingsPage(l)
	pages[WalletSettings] = WalletSettingsPage(l)
	pages[SecurityTools] = SecurityToolsPage(l)
	pages[Proposals] = ProposalsPage(l)
	pages[ProposalDetails] = ProposalDetailsPage(l)
	pages[Debug] = DebugPage(l)
	pages[Log] = LogPage(l)
	pages[Statistics] = StatPage(l)
	pages[About] = AboutPage(l)
	pages[Help] = HelpPage(l)
	pages[UTXO] = UTXOPage(l)
	pages[AccountDetails] = AcctDetailsPage(l)
	pages[Privacy] = PrivacyPage(l)
	pages[Tickets] = TicketPage(l)
	pages[ValidateAddress] = ValidateAddressPage(l)
	pages[TicketsList] = TicketPageList(l)
	pages[TicketsActivity] = TicketActivityPage(l)

	return pages
}