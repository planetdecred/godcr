// SPDX-License-Identifier: Unlicense OR MIT

package materialplus

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"gioui.org/f32"
	"gioui.org/gesture"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"

	"golang.org/x/image/math/fixed"
)

// Editor implements an editable and scrollable text area.
type Editor struct {
	Font text.Font
	// Color is the text color.
	Color color.RGBA
	// Hint contains the text displayed when the editor is empty.
	Hint string
	// HintColor is the color of hint text.
	HintColor color.RGBA

	shaper text.Shaper

	Alignment text.Alignment
	// SingleLine force the text to stay on a single line.
	// SingleLine also sets the scrolling direction to
	// horizontal.
	SingleLine bool
	// Submit enabled translation of carriage return keys to SubmitEvents.
	// If not enabled, carriage returns are inserted as newlines in the text.
	Submit bool

	// IsNumeric forces the editor to only accept numeric characters if set to true.
	isNumeric bool

	eventKey     int
	scale        int
	font         text.Font
	blinkStart   time.Time
	focused      bool
	rr           editBuffer
	mask         string
	maxWidth     int
	viewSize     image.Point
	valid        bool
	lines        []text.Line
	shapes       []line
	dims         layout.Dimensions
	carWidth     fixed.Int26_6
	requestFocus bool
	caretOn      bool
	caretScroll  bool

	// carXOff is the offset to the current caret
	// position when moving between lines.
	carXOff fixed.Int26_6

	scroller  gesture.Scroll
	scrollOff image.Point

	clicker gesture.Click

	// events is the list of events not yet processed.
	events []EditorEvent
	// prevEvents is the number of events from the previous frame.
	prevEvents int
}

// EditorEvent represents an event in the liftime of the editor
type EditorEvent interface {
	isEditorEvent()
}

// A ChangeEvent is generated for every user change to the text.
type ChangeEvent struct{}

// A SubmitEvent is generated when Submit is set
// and a carriage return key is pressed.
type SubmitEvent struct {
	Text string
}

type line struct {
	offset f32.Point
	clip   op.CallOp
}

const (
	inf              = 1e6
	blinksPerSecond  = 1
	maxBlinkDuration = 10 * time.Second
)

// Editor returns an instance of Editor
func (t *Theme) Editor(hint, mask string) *Editor {
	return &Editor{
		Font: text.Font{
			Size: t.TextSize,
		},
		Color:     t.Color.Text,
		shaper:    t.Shaper,
		Hint:      hint,
		mask:      mask,
		HintColor: t.Color.Hint,
	}
}

// IsNumeric sets the editor to only accept numeric characters
func (e *Editor) IsNumeric() {
	e.isNumeric = true
}

// Events returns available editor events.
func (e *Editor) Events(gtx *layout.Context) []EditorEvent {
	e.processEvents(gtx)
	events := e.events
	e.events = nil
	e.prevEvents = 0
	return events
}

func (e *Editor) processEvents(gtx *layout.Context) {
	e.processPointer(gtx)
	e.processKey(gtx)
}

func (e *Editor) processPointer(gtx *layout.Context) {
	sbounds := e.scrollBounds()
	var smin, smax int
	var axis gesture.Axis
	if e.SingleLine {
		axis = gesture.Horizontal
		smin, smax = sbounds.Min.X, sbounds.Max.X
	} else {
		axis = gesture.Vertical
		smin, smax = sbounds.Min.Y, sbounds.Max.Y
	}
	sdist := e.scroller.Scroll(gtx, gtx, gtx.Now(), axis)
	var soff int
	if e.SingleLine {
		e.scrollRel(sdist, 0)
		soff = e.scrollOff.X
	} else {
		e.scrollRel(0, sdist)
		soff = e.scrollOff.Y
	}
	for _, evt := range e.clicker.Events(gtx) {
		switch {
		case evt.Type == gesture.TypePress && evt.Source == pointer.Mouse,
			evt.Type == gesture.TypeClick && evt.Source == pointer.Touch:
			e.blinkStart = gtx.Now()
			e.moveCoord(gtx, image.Point{
				X: int(math.Round(float64(evt.Position.X))),
				Y: int(math.Round(float64(evt.Position.Y))),
			})
			e.requestFocus = true
			if e.scroller.State() != gesture.StateFlinging {
				e.caretScroll = true
			}
		}
	}
	if (sdist > 0 && soff >= smax) || (sdist < 0 && soff <= smin) {
		e.scroller.Stop()
	}
}

func (e *Editor) processKey(gtx *layout.Context) {
	for _, ke := range gtx.Events(&e.eventKey) {
		e.blinkStart = gtx.Now()
		switch ke := ke.(type) {
		case key.FocusEvent:
			e.focused = ke.Focus
		case key.Event:
			if !e.focused {
				break
			}
			if e.Submit && (ke.Name == key.NameReturn || ke.Name == key.NameEnter) {
				if !ke.Modifiers.Contain(key.ModShift) {
					e.events = append(e.events, SubmitEvent{
						Text: e.Text(),
					})
					return
				}
			}
			if e.command(ke) {
				e.caretScroll = true
				e.scroller.Stop()
			}
		case key.EditEvent:
			e.caretScroll = true
			e.scroller.Stop()
			e.append(ke.Text)
		}
		if e.rr.Changed() {
			e.events = append(e.events, ChangeEvent{})
		}
	}
}

func (e *Editor) command(k key.Event) bool {
	switch k.Name {
	case key.NameReturn, key.NameEnter:
		e.append("\n")
	case key.NameDeleteBackward:
		e.Delete(-1)
	case key.NameDeleteForward:
		e.Delete(1)
	case key.NameUpArrow:
		line, _, carX, _ := e.layoutCaret()
		e.carXOff = e.moveToLine(carX+e.carXOff, line-1)
	case key.NameDownArrow:
		line, _, carX, _ := e.layoutCaret()
		e.carXOff = e.moveToLine(carX+e.carXOff, line+1)
	case key.NameLeftArrow:
		e.Move(-1)
	case key.NameRightArrow:
		e.Move(1)
	case key.NamePageUp:
		e.movePages(-1)
	case key.NamePageDown:
		e.movePages(+1)
	case key.NameHome:
		e.moveStart()
	case key.NameEnd:
		e.moveEnd()
	default:
		return false
	}
	return true
}

// Focus requests the input focus for the Editor.
func (e *Editor) Focus() {
	e.requestFocus = true
}

// Draw draws the editor
func (e *Editor) Draw(gtx *layout.Context) {
	var stack op.StackOp
	stack.Push(gtx.Ops)
	var macro op.MacroOp
	macro.Record(gtx.Ops)
	paint.ColorOp{Color: e.HintColor}.Add(gtx.Ops)
	tl := widget.Label{Alignment: e.Alignment}
	tl.Layout(gtx, e.shaper, e.Font, e.Hint)
	macro.Stop()
	if w := gtx.Dimensions.Size.X; gtx.Constraints.Width.Min < w {
		gtx.Constraints.Width.Min = w
	}
	if h := gtx.Dimensions.Size.Y; gtx.Constraints.Height.Min < h {
		gtx.Constraints.Height.Min = h
	}
	e.Layout(gtx, e.shaper, e.Font)
	if e.Len() > 0 {
		paint.ColorOp{Color: e.Color}.Add(gtx.Ops)
		e.PaintText(gtx)
	} else {
		macro.Add()
	}
	paint.ColorOp{Color: e.Color}.Add(gtx.Ops)
	e.PaintCaret(gtx)
	stack.Pop()
}

// Layout lays out the editor.
func (e *Editor) Layout(gtx *layout.Context, sh text.Shaper, font text.Font) {
	// Flush events from before the previous frame.
	copy(e.events, e.events[e.prevEvents:])
	e.events = e.events[:len(e.events)-e.prevEvents]
	e.prevEvents = len(e.events)
	if e.font != font {
		e.invalidate()
		e.font = font
	}
	e.processEvents(gtx)
	e.layout(gtx, sh)
}

func (e *Editor) layout(gtx *layout.Context, sh text.Shaper) {
	// Crude configuration change detection.
	if scale := gtx.Px(unit.Sp(100)); scale != e.scale {
		e.invalidate()
		e.scale = scale
	}
	cs := gtx.Constraints
	e.carWidth = fixed.I(gtx.Px(unit.Dp(1)))

	maxWidth := cs.Width.Max
	if e.SingleLine {
		maxWidth = inf
	}
	if maxWidth != e.maxWidth {
		e.maxWidth = maxWidth
		e.invalidate()
	}

	if !e.valid {
		e.lines, e.dims = e.layoutText(gtx, sh, e.font)
		e.valid = true
	}

	e.viewSize = cs.Constrain(e.dims.Size)
	// Adjust scrolling for new viewport and layout.
	e.scrollRel(0, 0)

	if e.caretScroll {
		e.caretScroll = false
		e.scrollToCaret()
	}

	off := image.Point{
		X: -e.scrollOff.X,
		Y: -e.scrollOff.Y,
	}
	clip := textPadding(e.lines)
	clip.Max = clip.Max.Add(e.viewSize)
	it := lineIterator{
		Lines:     e.lines,
		Clip:      clip,
		Alignment: e.Alignment,
		Width:     e.viewSize.X,
		Offset:    off,
	}
	e.shapes = e.shapes[:0]
	for {
		_, _, layout, off, ok := it.Next()
		if !ok {
			break
		}
		path := sh.Shape(gtx, e.font, layout)
		e.shapes = append(e.shapes, line{off, path})
	}

	key.InputOp{Key: &e.eventKey, Focus: e.requestFocus}.Add(gtx.Ops)
	e.requestFocus = false
	pointerPadding := gtx.Px(unit.Dp(4))
	r := image.Rectangle{Max: e.viewSize}
	r.Min.X -= pointerPadding
	r.Min.Y -= pointerPadding
	r.Max.X += pointerPadding
	r.Max.X += pointerPadding
	pointer.Rect(r).Add(gtx.Ops)
	e.scroller.Add(gtx.Ops)
	e.clicker.Add(gtx.Ops)
	e.caretOn = false
	if e.focused {
		now := gtx.Now()
		dt := now.Sub(e.blinkStart)
		blinking := dt < maxBlinkDuration
		const timePerBlink = time.Second / blinksPerSecond
		nextBlink := now.Add(timePerBlink/2 - dt%(timePerBlink/2))
		if blinking {
			redraw := op.InvalidateOp{At: nextBlink}
			redraw.Add(gtx.Ops)
		}
		e.caretOn = e.focused && (!blinking || dt%timePerBlink < timePerBlink/2)
	}

	gtx.Dimensions = layout.Dimensions{Size: e.viewSize, Baseline: e.dims.Baseline}
}

// PaintText paints the text
func (e *Editor) PaintText(gtx *layout.Context) {
	clip := textPadding(e.lines)
	clip.Max = clip.Max.Add(e.viewSize)
	for _, shape := range e.shapes {
		var stack op.StackOp
		stack.Push(gtx.Ops)
		op.TransformOp{}.Offset(shape.offset).Add(gtx.Ops)
		shape.clip.Add(gtx.Ops)
		paint.PaintOp{Rect: toRectF(clip).Sub(shape.offset)}.Add(gtx.Ops)
		stack.Pop()
	}
}

// PaintCaret draws the blinker
func (e *Editor) PaintCaret(gtx *layout.Context) {
	if !e.caretOn {
		return
	}
	carLine, _, carX, carY := e.layoutCaret()

	var stack op.StackOp
	stack.Push(gtx.Ops)
	carX -= e.carWidth / 2
	carAsc, carDesc := -e.lines[carLine].Bounds.Min.Y, e.lines[carLine].Bounds.Max.Y
	carRect := image.Rectangle{
		Min: image.Point{X: carX.Ceil(), Y: carY - carAsc.Ceil()},
		Max: image.Point{X: carX.Ceil() + e.carWidth.Ceil(), Y: carY + carDesc.Ceil()},
	}
	carRect = carRect.Add(image.Point{
		X: -e.scrollOff.X,
		Y: -e.scrollOff.Y,
	})
	clip := textPadding(e.lines)
	// Account for caret width to each side.
	whalf := (e.carWidth / 2).Ceil()
	if clip.Max.X < whalf {
		clip.Max.X = whalf
	}
	if clip.Min.X > -whalf {
		clip.Min.X = -whalf
	}
	clip.Max = clip.Max.Add(e.viewSize)
	carRect = clip.Intersect(carRect)
	if !carRect.Empty() {
		paint.PaintOp{Rect: toRectF(carRect)}.Add(gtx.Ops)
	}
	stack.Pop()
}

// Len is the length of the editor contents.
func (e *Editor) Len() int {
	return e.rr.len()
}

// Text returns the contents of the editor.
func (e *Editor) Text() string {
	return e.rr.String()
}

// SetText replaces the contents of the editor.
func (e *Editor) SetText(s string) {
	// return if it's a numeric input but a string is entered
	if e.isNumeric && !validateNumberStr(s) {
		return
	}

	e.rr = editBuffer{}
	e.carXOff = 0
	e.prepend(s)
}

func (e *Editor) scrollBounds() image.Rectangle {
	var b image.Rectangle
	if e.SingleLine {
		if len(e.lines) > 0 {
			b.Min.X = align(e.Alignment, e.lines[0].Width, e.viewSize.X).Floor()
			if b.Min.X > 0 {
				b.Min.X = 0
			}
		}
		b.Max.X = e.dims.Size.X + b.Min.X - e.viewSize.X
	} else {
		b.Max.Y = e.dims.Size.Y - e.viewSize.Y
	}
	return b
}

func (e *Editor) scrollRel(dx, dy int) {
	e.scrollAbs(e.scrollOff.X+dx, e.scrollOff.Y+dy)
}

func (e *Editor) scrollAbs(x, y int) {
	e.scrollOff.X = x
	e.scrollOff.Y = y
	b := e.scrollBounds()
	if e.scrollOff.X > b.Max.X {
		e.scrollOff.X = b.Max.X
	}
	if e.scrollOff.X < b.Min.X {
		e.scrollOff.X = b.Min.X
	}
	if e.scrollOff.Y > b.Max.Y {
		e.scrollOff.Y = b.Max.Y
	}
	if e.scrollOff.Y < b.Min.Y {
		e.scrollOff.Y = b.Min.Y
	}
}

func (e *Editor) moveCoord(c unit.Converter, pos image.Point) {
	var (
		prevDesc fixed.Int26_6
		carLine  int
		y        int
	)
	for _, l := range e.lines {
		y += (prevDesc + l.Ascent).Ceil()
		prevDesc = l.Descent
		if y+prevDesc.Ceil() >= pos.Y+e.scrollOff.Y {
			break
		}
		carLine++
	}
	x := fixed.I(pos.X + e.scrollOff.X)
	e.moveToLine(x, carLine)
}

func (e *Editor) layoutText(c unit.Converter, s text.Shaper, font text.Font) ([]text.Line, layout.Dimensions) {
	e.rr.Reset()

	buff := e.rr
	if e.mask != "" {
		buff = editBuffer{}
		buff.prepend(strings.Repeat(e.mask, e.Len()))
	}

	opts := text.LayoutOptions{MaxWidth: e.maxWidth}
	lines, _ := s.Layout(c, font, &buff, opts)
	dims := linesDimens(lines)
	for i := 0; i < len(lines)-1; i++ {
		// To avoid layout flickering while editing, assume a soft newline takes
		// up all available space.
		if layout := lines[i].Layout; len(layout) > 0 {
			r := layout[len(layout)-1].Rune
			if r != '\n' {
				dims.Size.X = e.maxWidth
				break
			}
		}
	}
	return lines, dims
}

func (e *Editor) layoutCaret() (carLine, carCol int, x fixed.Int26_6, y int) {
	var idx int
	var prevDesc fixed.Int26_6
loop:
	for carLine = 0; carLine < len(e.lines); carLine++ {
		l := e.lines[carLine]
		y += (prevDesc + l.Ascent).Ceil()
		prevDesc = l.Descent
		if carLine == len(e.lines)-1 || idx+len(l.Layout) > e.rr.caret {
			for _, g := range l.Layout {
				if idx == e.rr.caret {
					break loop
				}
				x += g.Advance
				idx += utf8.RuneLen(g.Rune)
				carCol++
			}
			break
		}
		idx += l.Len
	}
	x += align(e.Alignment, e.lines[carLine].Width, e.viewSize.X)
	return
}

func (e *Editor) invalidate() {
	e.valid = false
}

// Delete runes from the caret position. The sign of runes specifies the
// direction to delete: positive is forward, negative is backward.
func (e *Editor) Delete(runes int) {
	e.rr.deleteRunes(runes)
	e.carXOff = 0
	e.invalidate()
}

// Insert inserts text at the caret, moving the caret forward.
func (e *Editor) Insert(s string) {
	// return if it's a numeric input but a string is entered
	if e.isNumeric && !validateNumberStr(s) {
		return
	}

	e.append(s)
	e.caretScroll = true
	e.invalidate()
}

func (e *Editor) append(s string) {
	// return if it's a numeric input but a string is entered
	if e.isNumeric && !validateNumberStr(s) {
		return
	}

	if e.SingleLine {
		s = strings.ReplaceAll(s, "\n", "")
	}
	e.prepend(s)
	e.rr.caret += len(s)
}

func (e *Editor) prepend(s string) {
	e.rr.prepend(s)
	e.carXOff = 0
	e.invalidate()
}

func (e *Editor) movePages(pages int) {
	_, _, carX, carY := e.layoutCaret()
	y := carY + pages*e.viewSize.Y
	var (
		prevDesc fixed.Int26_6
		carLine2 int
	)
	y2 := e.lines[0].Ascent.Ceil()
	for i := 1; i < len(e.lines); i++ {
		if y2 >= y {
			break
		}
		l := e.lines[i]
		h := (prevDesc + l.Ascent).Ceil()
		prevDesc = l.Descent
		if y2+h-y >= y-y2 {
			break
		}
		y2 += h
		carLine2++
	}
	e.carXOff = e.moveToLine(carX+e.carXOff, carLine2)
}

func (e *Editor) moveToLine(carX fixed.Int26_6, carLine2 int) fixed.Int26_6 {
	carLine, carCol, _, _ := e.layoutCaret()
	if carLine2 < 0 {
		carLine2 = 0
	}
	if carLine2 >= len(e.lines) {
		carLine2 = len(e.lines) - 1
	}
	// Move to start of line.
	for i := carCol - 1; i >= 0; i-- {
		_, s := e.rr.runeBefore(e.rr.caret)
		e.rr.caret -= s
	}
	if carLine2 != carLine {
		// Move to start of line2.
		if carLine2 > carLine {
			for i := carLine; i < carLine2; i++ {
				e.rr.caret += e.lines[i].Len
			}
		} else {
			for i := carLine - 1; i >= carLine2; i-- {
				e.rr.caret -= e.lines[i].Len
			}
		}
	}
	l2 := e.lines[carLine2]
	carX2 := align(e.Alignment, l2.Width, e.viewSize.X)
	// Only move past the end of the last line
	end := 0
	if carLine2 < len(e.lines)-1 {
		end = 1
	}
	// Move to rune closest to previous horizontal position.
	for i := 0; i < len(l2.Layout)-end; i++ {
		g := l2.Layout[i]
		if carX2 >= carX {
			break
		}
		if carX2+g.Advance-carX >= carX-carX2 {
			break
		}
		carX2 += g.Advance
		_, s := e.rr.runeAt(e.rr.caret)
		e.rr.caret += s
	}
	return carX - carX2
}

// Move the caret: positive distance moves forward, negative distance moves
// backward.
func (e *Editor) Move(distance int) {
	e.rr.move(distance)
	e.carXOff = 0
}

func (e *Editor) moveStart() {
	carLine, carCol, x, _ := e.layoutCaret()
	layout := e.lines[carLine].Layout
	for i := carCol - 1; i >= 0; i-- {
		_, s := e.rr.runeBefore(e.rr.caret)
		e.rr.caret -= s
		x -= layout[i].Advance
	}
	e.carXOff = -x
}

func (e *Editor) moveEnd() {
	carLine, carCol, x, _ := e.layoutCaret()
	l := e.lines[carLine]
	// Only move past the end of the last line
	end := 0
	if carLine < len(e.lines)-1 {
		end = 1
	}
	layout := l.Layout
	for i := carCol; i < len(layout)-end; i++ {
		adv := layout[i].Advance
		_, s := e.rr.runeAt(e.rr.caret)
		e.rr.caret += s
		x += adv
	}
	a := align(e.Alignment, l.Width, e.viewSize.X)
	e.carXOff = l.Width + a - x
}

func (e *Editor) scrollToCaret() {
	carLine, _, x, y := e.layoutCaret()
	l := e.lines[carLine]
	if e.SingleLine {
		var dist int
		if d := x.Floor() - e.scrollOff.X; d < 0 {
			dist = d
		} else if d := x.Ceil() - (e.scrollOff.X + e.viewSize.X); d > 0 {
			dist = d
		}
		e.scrollRel(dist, 0)
	} else {
		miny := y - l.Ascent.Ceil()
		maxy := y + l.Descent.Ceil()
		var dist int
		if d := miny - e.scrollOff.Y; d < 0 {
			dist = d
		} else if d := maxy - (e.scrollOff.Y + e.viewSize.Y); d > 0 {
			dist = d
		}
		e.scrollRel(0, dist)
	}
}

func (s ChangeEvent) isEditorEvent() {}
func (s SubmitEvent) isEditorEvent() {}

type lineIterator struct {
	Lines     []text.Line
	Clip      image.Rectangle
	Alignment text.Alignment
	Width     int
	Offset    image.Point

	y, prevDesc fixed.Int26_6
	txtOff      int
}

func (l *lineIterator) Next() (int, int, []text.Glyph, f32.Point, bool) {
	for len(l.Lines) > 0 {
		line := l.Lines[0]
		l.Lines = l.Lines[1:]
		x := align(l.Alignment, line.Width, l.Width) + fixed.I(l.Offset.X)
		l.y += l.prevDesc + line.Ascent
		l.prevDesc = line.Descent
		// Align baseline and line start to the pixel grid.
		off := fixed.Point26_6{X: fixed.I(x.Floor()), Y: fixed.I(l.y.Ceil())}
		l.y = off.Y
		off.Y += fixed.I(l.Offset.Y)
		if (off.Y + line.Bounds.Min.Y).Floor() > l.Clip.Max.Y {
			break
		}
		layout := line.Layout
		start := l.txtOff
		l.txtOff += line.Len
		if (off.Y + line.Bounds.Max.Y).Ceil() < l.Clip.Min.Y {
			continue
		}
		for len(layout) > 0 {
			g := layout[0]
			adv := g.Advance
			if (off.X + adv + line.Bounds.Max.X - line.Width).Ceil() >= l.Clip.Min.X {
				break
			}
			off.X += adv
			layout = layout[1:]
			start += utf8.RuneLen(g.Rune)
		}
		end := start
		endx := off.X
		for i, g := range layout {
			if (endx + line.Bounds.Min.X).Floor() > l.Clip.Max.X {
				layout = layout[:i]
				break
			}
			end += utf8.RuneLen(g.Rune)
			endx += g.Advance
		}
		offf := f32.Point{X: float32(off.X) / 64, Y: float32(off.Y) / 64}
		return start, end, layout, offf, true
	}
	return 0, 0, nil, f32.Point{}, false
}

func toRectF(r image.Rectangle) f32.Rectangle {
	return f32.Rectangle{
		Min: f32.Point{X: float32(r.Min.X), Y: float32(r.Min.Y)},
		Max: f32.Point{X: float32(r.Max.X), Y: float32(r.Max.Y)},
	}
}

func textPadding(lines []text.Line) (padding image.Rectangle) {
	if len(lines) == 0 {
		return
	}
	first := lines[0]
	if d := first.Ascent + first.Bounds.Min.Y; d < 0 {
		padding.Min.Y = d.Ceil()
	}
	last := lines[len(lines)-1]
	if d := last.Bounds.Max.Y - last.Descent; d > 0 {
		padding.Max.Y = d.Ceil()
	}
	if d := first.Bounds.Min.X; d < 0 {
		padding.Min.X = d.Ceil()
	}
	if d := first.Bounds.Max.X - first.Width; d > 0 {
		padding.Max.X = d.Ceil()
	}
	return
}

func linesDimens(lines []text.Line) layout.Dimensions {
	var width fixed.Int26_6
	var h int
	var baseline int
	if len(lines) > 0 {
		baseline = lines[0].Ascent.Ceil()
		var prevDesc fixed.Int26_6
		for _, l := range lines {
			h += (prevDesc + l.Ascent).Ceil()
			prevDesc = l.Descent
			if l.Width > width {
				width = l.Width
			}
		}
		h += lines[len(lines)-1].Descent.Ceil()
	}
	w := width.Ceil()
	return layout.Dimensions{
		Size: image.Point{
			X: w,
			Y: h,
		},
		Baseline: h - baseline,
	}
}

func align(align text.Alignment, width fixed.Int26_6, maxWidth int) fixed.Int26_6 {
	mw := fixed.I(maxWidth)
	switch align {
	case text.Middle:
		return fixed.I(((mw - width) / 2).Floor())
	case text.End:
		return fixed.I((mw - width).Floor())
	case text.Start:
		return 0
	default:
		panic(fmt.Errorf("unknown alignment %v", align))
	}
}

func validateNumberStr(s string) bool {
	if _, err := strconv.Atoi(s); err != nil {
		return false
	}

	return true
}
