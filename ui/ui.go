package ui

import (
	"fmt"
	"image"
	"os"
	"path"
	"path/filepath"
	"strings"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/font/opentype"
	"gioui.org/io/clipboard"
	"gioui.org/io/key"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/text"

	"github.com/planetdecred/godcr/dexc"
	"github.com/planetdecred/godcr/ui/uidex"
	"github.com/planetdecred/godcr/ui/uiwallet"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/planetdecred/godcr/wallet"
)

// UI represents the wallet UI and DEX UI of the app
type UI struct {
	view     int
	uiwallet *uiwallet.Wallet
	uidex    *uidex.DexUI
}

// NewUI creates and initializes a new ui with start
// as the first page displayed.
// Should never be called more than once as it calls
func NewUI(w *app.Window, wal *wallet.Wallet, dexc *dexc.Dexc, internalLog chan string) (*UI, error) {
	ui := &UI{
		view: 1,
	}

	absoluteWdPath, err := getAbsoultePath()
	if err != nil {
		panic(err)
	}

	// Initialize font face
	var collection []text.FontFace
	source, err := os.Open(filepath.Join(absoluteWdPath, "ui/assets/fonts/source_sans_pro_regular.otf"))
	if err != nil {
		fmt.Println("Failed to load font")
		collection = gofont.Collection()
	} else {
		stat, err := source.Stat()
		if err != nil {
			fmt.Println(err)
		}
		bytes := make([]byte, stat.Size())
		source.Read(bytes)
		fnt, err := opentype.Parse(bytes)
		if err != nil {
			fmt.Println(err)
		}
		collection = append(collection, text.FontFace{Font: text.Font{}, Face: fnt})
	}

	// Initialize wallet icons
	walletIcons, err := initIcons("ui/assets/walleticons")
	if err != nil {
		return nil, err
	}

	// Initialize wallet icons
	dexIcons, err := initIcons("ui/assets/dexicons")
	if err != nil {
		return nil, err
	}

	// Create wallet ui
	uiw, err := uiwallet.NewWalletUI(wal, walletIcons, collection, internalLog, &ui.view, w.Invalidate)
	if err != nil {
		return nil, fmt.Errorf("Could not initialize wallet UI: %s\ns", err)
	}

	// Create Dex ui
	uid, err := uidex.NewDexUI(dexc, dexIcons, collection, internalLog, &ui.view, w.Invalidate)
	if err != nil {
		return nil, fmt.Errorf("Could not initialize dex UI: %s\ns", err)
	}

	ui.uidex = uid
	ui.uiwallet = uiw

	return ui, nil
}

func (ui *UI) Loop(shutdown chan int, w *app.Window) error {
	go func() {
		for e := range w.Events() {
			switch e := e.(type) {
			case system.DestroyEvent:
				ui.uiwallet.HandlerDestroy(shutdown)
			case system.FrameEvent:
				var gtx layout.Context
				if ui.view == values.WalletView {
					gtx = layout.NewContext(ui.uiwallet.Ops(), e)
					ui.uiwallet.HandlerPages(gtx)
				} else {
					gtx = layout.NewContext(ui.uidex.Ops(), e)
					ui.uidex.HandlerPages(gtx)
				}

				e.Frame(gtx.Ops)
			case key.Event:
				ui.uiwallet.HandlerKeyEvents(&e)
			case clipboard.Event:
				ui.uiwallet.HandlerClipboard(&e)
			case nil:
				// Ignore
			default:
				// log.Tracef("Unhandled window event %+v\n", e)
			}
		}
	}()

	go ui.uiwallet.Run(shutdown, w)
	go ui.uidex.Run(shutdown, w)
	return nil
}

func getAbsoultePath() (string, error) {
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

func initIcons(relativePath string) (map[string]image.Image, error) {
	icons := make(map[string]image.Image)
	absoluteWdPath, err := getAbsoultePath()
	if err != nil {
		panic(err)
	}
	err = filepath.Walk(filepath.Join(absoluteWdPath, relativePath), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			panic(err)
		}
		if info.IsDir() || !strings.HasSuffix(path, ".png") {
			return nil
		}

		f, _ := os.Open(path)
		img, _, err := image.Decode(f)
		if err != nil {
			return err
		}
		split := strings.Split(info.Name(), ".")
		icons[split[0]] = img
		return nil
	})

	if err != nil {
		return nil, err
	}

	return icons, nil
}
