package icons

import (
	"fmt"

	"github.com/markbates/pkger"
)

const (
	Receive = "/ui/icons/icons.go"
)

func readFile(path string) (blob []byte, err error) {
	source, err := pkger.Open(path)
	if err != nil {
		return
	}
	_, err = source.Read(blob)
	if err != nil {
		return
	}

	return
}

func RenderIcon(icon string, size int) {
	logo, err := readFile(icon)
	if err != nil {
		fmt.Printf("RENDER ICON ERROR %v\n \n", err.Error())
		return
	}
	fmt.Printf("Logo %v", logo)
	// m, _:= iconvg.DecodeMetadata(logo)
	//dx, dy := m.ViewBox.AspectRatio()
	//img := image.NewRGBA(image.Rectangle{Max: image.Point{X: size, Y: int(float32(size) * dy / dx)}})
	//var ico iconvg.Rasterizer
	//ico.SetDstImage(img, img.Bounds(), draw.Src)
	//// Use white for icons.
	//m.Palette[0] = color.RGBA{A: 0xff, R: 0x10, G: 0x10, B: 0x10}
	//err = iconvg.Decode(&ico, logo, &iconvg.DecodeOptions{
	//	Palette: &m.Palette,
	//})
	//if err != nil {
	//	fmt.Printf("RENDER ICON ERROR %v", err.Error())
	//	return
	//}
	//return
}
