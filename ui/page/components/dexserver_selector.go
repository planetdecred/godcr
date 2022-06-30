package components

import (
	"sort"
	"strings"

	"decred.org/dcrdex/client/core"
	"gioui.org/layout"
	"gioui.org/widget/material"
	"github.com/planetdecred/godcr/app"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/values"
)

const DexServerSelectorID = "dex_server_selector"

type DexServerSelector struct {
	*load.Load
	// GenericPageModal defines methods such as ID() and OnAttachedToNavigator()
	// that helps this Page satisfy the app.Page interface. It also defines
	// helper methods for accessing the PageNavigator that displayed this page
	// and the root WindowNavigator.
	*app.GenericPageModal

	shadowBox       *decredmaterial.Shadow
	knownDexServers *decredmaterial.ClickableList
	materialLoader  material.LoaderStyle

	dexServerSelected func(server string)
}

func NewDexServerSelector(l *load.Load, onDexServerSelected func(server string)) *DexServerSelector {
	ds := &DexServerSelector{
		GenericPageModal:  app.NewGenericPageModal(DexServerSelectorID),
		Load:              l,
		dexServerSelected: onDexServerSelected,
		shadowBox:         l.Theme.Shadow(),
		materialLoader:    material.Loader(l.Theme.Base),
	}

	ds.knownDexServers = l.Theme.NewClickableList(layout.Vertical)

	return ds
}

func (ds *DexServerSelector) Expose() {
	ds.startDexClient()
}

// isLoadingDexClient check for Dexc start, initialized, loggedin status,
// since Dex client UI not required for app password, IsInitialized and IsLoggedIn should be done at dcrlibwallet.
func (ds *DexServerSelector) isLoadingDexClient() bool {
	return ds.Dexc().Core() == nil || !ds.Dexc().Core().IsInitialized() || !ds.Dexc().IsLoggedIn()
}

// startDexClient do start DEX client,
// initialize and login to DEX,
// since Dex client UI not required for app password, initialize and login should be done at dcrlibwallet.
func (ds *DexServerSelector) startDexClient() {
	_, err := ds.WL.MultiWallet.StartDexClient()
	if err != nil {
		ds.Toast.NotifyError(err.Error())
		return
	}

	// TODO: move to dcrlibwallet sine bypass Dex password by DEXClientPass
	if !ds.Dexc().Initialized() {
		err = ds.Dexc().InitializeWithPassword([]byte(values.DEXClientPass))
		if err != nil {
			ds.Toast.NotifyError(err.Error())
			return
		}
	}

	if !ds.Dexc().IsLoggedIn() {
		err := ds.Dexc().Login([]byte(values.DEXClientPass))
		if err != nil {
			ds.Toast.NotifyError(err.Error())
			return
		}
	}
}

func (ds *DexServerSelector) DexServersLayout(gtx C) D {
	if ds.isLoadingDexClient() {
		return layout.Center.Layout(gtx, func(gtx C) D {
			gtx.Constraints.Min.X = 50
			return ds.materialLoader.Layout(gtx)
		})
	}
	knownDexServers := ds.mapKnowDexServers()
	dexServers := sortDexExchanges(knownDexServers)
	return ds.knownDexServers.Layout(gtx, len(dexServers), func(gtx C, i int) D {
		dexServer := dexServers[i]
		hostport := strings.Split(dexServer, ":")
		host := hostport[0]
		ds.shadowBox.SetShadowRadius(14)

		return decredmaterial.LinearLayout{
			Width:      decredmaterial.WrapContent,
			Height:     decredmaterial.WrapContent,
			Padding:    layout.UniformInset(values.MarginPadding18),
			Background: ds.Theme.Color.Surface,
			Alignment:  layout.Middle,
			Shadow:     ds.shadowBox,
			Margin:     layout.UniformInset(values.MarginPadding5),
			Border:     decredmaterial.Border{Radius: decredmaterial.Radius(14)},
		}.Layout(gtx,
			layout.Flexed(1, ds.Theme.Label(values.TextSize16, host).Layout),
		)
	})
}

// sortDexExchanges convert map cert into a sorted slice
func sortDexExchanges(mapCert map[string][]byte) []string {
	servers := make([]string, 0, len(mapCert))
	for host := range mapCert {
		servers = append(servers, host)
	}
	sort.Slice(servers, func(i, j int) bool {
		return servers[i] < servers[j]
	})
	return servers
}

func (ds *DexServerSelector) HandleUserInteractions() {
	if ok, index := ds.knownDexServers.ItemClicked(); ok {
		knownDexServers := ds.mapKnowDexServers()
		dexServers := sortDexExchanges(knownDexServers)
		ds.dexServerSelected(dexServers[index])
	}
}

// TODO: handler CRUD dex servers at dcrlibwallet
const KnownDexServersConfigKey = "known_dex_servers"

type DexServer struct {
	SavedHosts map[string][]byte
}

func (ds *DexServerSelector) mapKnowDexServers() map[string][]byte {
	knownDexServers := core.CertStore[ds.Dexc().Core().Network()]
	dexServer := new(DexServer)
	err := ds.Load.WL.MultiWallet.ReadUserConfigValue(KnownDexServersConfigKey, &dexServer)
	if err != nil {
		return knownDexServers
	}
	if dexServer.SavedHosts == nil {
		return knownDexServers
	}
	for host, cert := range dexServer.SavedHosts {
		knownDexServers[host] = cert
	}
	return knownDexServers
}
