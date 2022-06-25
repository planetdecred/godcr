package decredmaterial

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"math/rand"
	"strings"
	"time"

	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/widget"
)

func drawInk(gtx layout.Context, c widget.Press, highlightColor color.NRGBA) {
	// duration is the number of seconds for the
	// completed animation: expand while fading in, then
	// out.
	const (
		expandDuration = float32(0.5)
		fadeDuration   = float32(0.9)
	)

	now := gtx.Now

	t := float32(now.Sub(c.Start).Seconds())

	end := c.End
	if end.IsZero() {
		// If the press hasn't ended, don't fade-out.
		end = now
	}

	endt := float32(end.Sub(c.Start).Seconds())

	// Compute the fade-in/out position in [0;1].
	var alphat float32
	{
		var haste float32
		if c.Cancelled {
			// If the press was cancelled before the inkwell
			// was fully faded in, fast forward the animation
			// to match the fade-out.
			if h := 0.5 - endt/fadeDuration; h > 0 {
				haste = h
			}
		}
		// Fade in.
		half1 := t/fadeDuration + haste
		if half1 > 0.5 {
			half1 = 0.5
		}

		// Fade out.
		half2 := float32(now.Sub(end).Seconds())
		half2 /= fadeDuration
		half2 += haste
		if half2 > 0.5 {
			// Too old.
			return
		}

		alphat = half1 + half2
	}

	// Compute the expand position in [0;1].
	sizet := t
	if c.Cancelled {
		// Freeze expansion of cancelled presses.
		sizet = endt
	}
	sizet /= expandDuration

	// Animate only ended presses, and presses that are fading in.
	if !c.End.IsZero() || sizet <= 1.0 {
		op.InvalidateOp{}.Add(gtx.Ops)
	}

	if sizet > 1.0 {
		sizet = 1.0
	}

	if alphat > .5 {
		// Start fadeout after half the animation.
		alphat = 1.0 - alphat
	}
	// Twice the speed to attain fully faded in at 0.5.
	t2 := alphat * 2
	// Beziér ease-in curve.
	alphaBezier := t2 * t2 * (3.0 - 2.0*t2)
	sizeBezier := sizet * sizet * (3.0 - 2.0*sizet)
	size := gtx.Constraints.Min.X
	if h := gtx.Constraints.Min.Y; h > size {
		size = h
	}
	// Cover the entire constraints min rectangle.
	size *= 2 * int(math.Sqrt(2))
	// Apply curve values to size and color.
	size *= int(sizeBezier)
	alpha := 0.7 * alphaBezier
	const col = 0.8
	ba, _ := byte(alpha*0xff), byte(col*0xff)
	rgba := mulAlpha(highlightColor, ba)
	ink := paint.ColorOp{Color: rgba}
	ink.Add(gtx.Ops)
	rr := size / 2
	defer op.Offset(c.Position.Add(image.Point{
		X: -rr,
		Y: -rr,
	})).Push(gtx.Ops).Pop()
	defer clip.UniformRRect(image.Rectangle{Max: image.Pt(size, size)}, rr).Push(gtx.Ops).Pop()
	paint.PaintOp{}.Add(gtx.Ops)
}

func GenerateRandomNumber() int {
	return rand.New(rand.NewSource(time.Now().UnixNano())).Int()
}

// AnyKey returns a key.Set that will match any of the provided keys.
// See the doc on key.Set for more.
func AnyKey(keys ...string) key.Set {
	keysCombined := strings.Join(keys, "|")
	return key.Set(keysCombined)
}

// AnyKeyWithOptionalModifier returns a key.Set that will match any of the
// provided keys whether or not they are pressed with provided modifier.
// See the doc on key.Set for more.
func AnyKeyWithOptionalModifier(modifier key.Modifiers, keys ...string) key.Set {
	keysCombined := strings.Join(keys, ",")
	keysWithModifier := fmt.Sprintf("(%s)-[%s]", modifier, keysCombined)
	return key.Set(keysWithModifier)
}
