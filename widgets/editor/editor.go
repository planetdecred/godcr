package editor

import (
	"image"
	"math"
	"strconv"
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

	"golang.org/x/image/math/fixed"

	"github.com/raedahgroup/godcr-gio/widgets"
)

// Editor implements an editable and scrollable text area.
type Editor struct {
	Alignment text.Alignment
	// SingleLine force the text to stay on a single line.
	// SingleLine also sets the scrolling direction to
	// horizontal.
	SingleLine bool
	// Submit enabled translation of carriage return keys to SubmitEvents.
	// If not enabled, carriage returns are inserted as newlines in the text.
	Submit bool

	eventKey     int
	scale        int
	font         text.Font
	blinkStart   time.Time
	focused      bool
	rr           editBuffer
	realText     editBuffer
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

	char      string
	isNumeric bool

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
	blinksPerSecond  = 1
	maxBlinkDuration = 10 * time.Second
	inf              = 1e6
)

func NewEditor() *Editor {
	return &Editor{}
}

func (e *Editor) setMask(char string) *Editor {
	e.char = char
	return e
}

func (e *Editor) numeric() *Editor {
	e.isNumeric = true
	return e
}

// Events returns available editor events.
func (e *Editor) Events(ctx *layout.Context) []EditorEvent {
	e.processEvents(ctx)
	events := e.events
	e.events = nil
	return events
}

func (e *Editor) processEvents(ctx *layout.Context) {
	e.processPointer(ctx)
	e.processKey(ctx)
}

func (e *Editor) processPointer(ctx *layout.Context) {
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
	sdist := e.scroller.Scroll(ctx, ctx, ctx.Now(), axis)
	var soff int
	if e.SingleLine {
		e.scrollRel(sdist, 0)
		soff = e.scrollOff.X
	} else {
		e.scrollRel(0, sdist)
		soff = e.scrollOff.Y
	}
	for _, evt := range e.clicker.Events(ctx) {
		switch {
		case evt.Type == gesture.TypePress && evt.Source == pointer.Mouse,
			evt.Type == gesture.TypeClick && evt.Source == pointer.Touch:
			e.blinkStart = ctx.Now()
			e.moveCoord(ctx, image.Point{
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

func (e *Editor) processKey(ctx *layout.Context) {
	for _, ke := range ctx.Events(&e.eventKey) {
		e.blinkStart = ctx.Now()
		switch ke := ke.(type) {
		case key.FocusEvent:
			e.focused = ke.Focus
		case key.Event:
			if !e.focused {
				break
			}
			if e.Submit && ke.Name == key.NameReturn || ke.Name == key.NameEnter {
				if !ke.Modifiers.Contain(key.ModShift) {
					e.events = append(e.events, SubmitEvent{
						Text: e.Text(),
					})
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
		if e.realText.Changed() {
			e.events = append(e.events, ChangeEvent{})
		}
	}
}

// Focus requests the input focus for the Editor.
func (e *Editor) Focus() {
	e.requestFocus = true
}

// Layout lays out the editor.
func (e *Editor) Layout(ctx *layout.Context, sh *text.Shaper, font text.Font) {
	// Flush events from before the previous frame.
	copy(e.events, e.events[e.prevEvents:])
	e.events = e.events[:len(e.events)-e.prevEvents]
	e.prevEvents = len(e.events)
	if e.font != font {
		e.invalidate()
		e.font = font
	}
	e.processEvents(ctx)
	e.layout(ctx, sh)
}

func (e *Editor) layout(ctx *layout.Context, sh *text.Shaper) {
	// Crude configuration change detection.
	if scale := ctx.Px(unit.Sp(100)); scale != e.scale {
		e.invalidate()
		e.scale = scale
	}
	cs := ctx.Constraints
	e.carWidth = fixed.I(ctx.Px(unit.Dp(1)))

	maxWidth := cs.Width.Max
	if e.SingleLine {
		maxWidth = inf
	}
	if maxWidth != e.maxWidth {
		e.maxWidth = maxWidth
		e.invalidate()
	}

	if !e.valid {
		e.lines, e.dims = e.layoutText(ctx, sh, e.font)
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
	clip := widgets.TextPadding(e.lines)
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
		str, off, ok := it.Next()
		if !ok {
			break
		}
		path := sh.Shape(ctx, e.font, str)
		e.shapes = append(e.shapes, line{off, path})
	}

	key.InputOp{Key: &e.eventKey, Focus: e.requestFocus}.Add(ctx.Ops)
	e.requestFocus = false
	pointerPadding := ctx.Px(unit.Dp(4))
	r := image.Rectangle{Max: e.viewSize}
	r.Min.X -= pointerPadding
	r.Min.Y -= pointerPadding
	r.Max.X += pointerPadding
	r.Max.X += pointerPadding
	pointer.Rect(r).Add(ctx.Ops)
	e.scroller.Add(ctx.Ops)
	e.clicker.Add(ctx.Ops)
	e.caretOn = false
	if e.focused {
		now := ctx.Now()
		dt := now.Sub(e.blinkStart)
		blinking := dt < maxBlinkDuration
		const timePerBlink = time.Second / blinksPerSecond
		nextBlink := now.Add(timePerBlink/2 - dt%(timePerBlink/2))
		if blinking {
			redraw := op.InvalidateOp{At: nextBlink}
			redraw.Add(ctx.Ops)
		}
		e.caretOn = e.focused && (!blinking || dt%timePerBlink < timePerBlink/2)
	}

	ctx.Dimensions = layout.Dimensions{Size: e.viewSize, Baseline: e.dims.Baseline}
}

func (e *Editor) PaintText(ctx *layout.Context) {
	clip := widgets.TextPadding(e.lines)
	clip.Max = clip.Max.Add(e.viewSize)
	for _, shape := range e.shapes {
		var stack op.StackOp
		stack.Push(ctx.Ops)
		op.TransformOp{}.Offset(shape.offset).Add(ctx.Ops)
		shape.clip.Add(ctx.Ops)
		paint.PaintOp{Rect: widgets.ToRectF(clip).Sub(shape.offset)}.Add(ctx.Ops)
		stack.Pop()
	}
}

func (e *Editor) PaintCaret(ctx *layout.Context) {
	if !e.caretOn {
		return
	}
	carLine, _, carX, carY := e.layoutCaret()

	var stack op.StackOp
	stack.Push(ctx.Ops)
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
	clip := widgets.TextPadding(e.lines)
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
		paint.PaintOp{Rect: widgets.ToRectF(carRect)}.Add(ctx.Ops)
	}
	stack.Pop()
}

// Len is the length of the editor contents.
func (e *Editor) Len() int {
	return e.rr.len()
}

// Text returns the contents of the editor.
func (e *Editor) Text() string {
	return e.realText.String()
}

// SetText replaces the contents of the editor.
func (e *Editor) SetText(s string) {
	e.rr, e.realText = editBuffer{}, editBuffer{}
	e.carXOff = 0

	// return if it's a numeric input but a string is entered
	if e.isNumeric {
		if _, err := strconv.Atoi(s); err != nil {
			return
		}
	}

	str := s
	if e.char != "" {
		str = e.char
	}
	e.prepend(str)
	e.prependRealText(s)
}

func (e *Editor) Clear() {
	e.rr, e.realText = editBuffer{}, editBuffer{}
}

func (e *Editor) scrollBounds() image.Rectangle {
	var b image.Rectangle
	if e.SingleLine {
		if len(e.lines) > 0 {
			b.Min.X = widgets.Align(e.Alignment, e.lines[0].Width, e.viewSize.X).Floor()
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

func (e *Editor) layoutText(c unit.Converter, s *text.Shaper, font text.Font) ([]text.Line, layout.Dimensions) {
	txt := e.rr.String()
	opts := text.LayoutOptions{MaxWidth: e.maxWidth}
	//opts := text.LayoutOptions{SingleLine: e.SingleLine, MaxWidth: e.maxWidth}
	textLayout := s.Layout(c, font, txt, opts)
	lines := textLayout.Lines
	dims := widgets.LinesDimens(lines)
	for i := 0; i < len(lines)-1; i++ {
		s := lines[i].Text.String
		// To avoid layout flickering while editing, assume a soft newline takes
		// up all available space.
		if len(s) > 0 {
			r, _ := utf8.DecodeLastRuneInString(s)
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
		if carLine == len(e.lines)-1 || idx+len(l.Text.String) > e.rr.caret {
			str := l.Text.String
			for _, adv := range l.Text.Advances {
				if idx == e.rr.caret {
					break loop
				}
				x += adv
				_, s := utf8.DecodeRuneInString(str)
				idx += s
				str = str[s:]
				carCol++
			}
			break
		}
		idx += len(l.Text.String)
	}
	x += widgets.Align(e.Alignment, e.lines[carLine].Width, e.viewSize.X)
	return
}

func (e *Editor) invalidate() {
	e.valid = false
}

func (e *Editor) deleteRune() {
	e.rr.deleteRune()
	e.realText.deleteRune()
	e.carXOff = 0
	e.invalidate()
}

func (e *Editor) deleteRuneForward() {
	e.rr.deleteRuneForward()
	e.realText.deleteRuneForward()
	e.carXOff = 0
	e.invalidate()
}

func (e *Editor) append(s string) {
	if e.SingleLine && s == "\n" {
		return
	}

	// return if it's a numeric input but a string is entered
	if e.isNumeric {
		if _, err := strconv.Atoi(s); err != nil {
			return
		}
	}

	/**str := s
	if e.char != "" {
		str = e.char
	}
	*/
	e.prepend(s)
	e.rr.caret += len(s)
	e.realText.caret += len(s)
}

func (e *Editor) prepend(s string) {
	str := s
	if e.char != "" {
		str = e.char
	}

	e.rr.prepend(str)
	e.prependRealText(s)
	e.carXOff = 0
	e.invalidate()
}

func (e *Editor) prependRealText(s string) {
	e.realText.prepend(s)
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
				e.rr.caret += len(e.lines[i].Text.String)
			}
		} else {
			for i := carLine - 1; i >= carLine2; i-- {
				e.rr.caret -= len(e.lines[i].Text.String)
			}
		}
	}
	l2 := e.lines[carLine2]
	carX2 := widgets.Align(e.Alignment, l2.Width, e.viewSize.X)
	// Only move past the end of the last line
	end := 0
	if carLine2 < len(e.lines)-1 {
		end = 1
	}
	// Move to rune closest to previous horizontal position.
	for i := 0; i < len(l2.Text.Advances)-end; i++ {
		adv := l2.Text.Advances[i]
		if carX2 >= carX {
			break
		}
		if carX2+adv-carX >= carX-carX2 {
			break
		}
		carX2 += adv
		_, s := e.rr.runeAt(e.rr.caret)
		e.rr.caret += s
	}
	return carX - carX2
}

func (e *Editor) moveLeft() {
	e.rr.moveLeft()
	e.carXOff = 0
}

func (e *Editor) moveRight() {
	e.rr.moveRight()
	e.carXOff = 0
}

func (e *Editor) moveStart() {
	carLine, carCol, x, _ := e.layoutCaret()
	advances := e.lines[carLine].Text.Advances
	for i := carCol - 1; i >= 0; i-- {
		_, s := e.rr.runeBefore(e.rr.caret)
		e.rr.caret -= s
		x -= advances[i]
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
	for i := carCol; i < len(l.Text.Advances)-end; i++ {
		adv := l.Text.Advances[i]
		_, s := e.rr.runeAt(e.rr.caret)
		e.rr.caret += s
		x += adv
	}
	a := widgets.Align(e.Alignment, l.Width, e.viewSize.X)
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

func (e *Editor) command(k key.Event) bool {
	switch k.Name {
	case key.NameReturn, key.NameEnter:
		e.append("\n")
	case key.NameDeleteBackward:
		e.deleteRune()
	case key.NameDeleteForward:
		e.deleteRuneForward()
	case key.NameUpArrow:
		line, _, carX, _ := e.layoutCaret()
		e.carXOff = e.moveToLine(carX+e.carXOff, line-1)
	case key.NameDownArrow:
		line, _, carX, _ := e.layoutCaret()
		e.carXOff = e.moveToLine(carX+e.carXOff, line+1)
	case key.NameLeftArrow:
		e.moveLeft()
	case key.NameRightArrow:
		e.moveRight()
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

func (s ChangeEvent) isEditorEvent() {}
func (s SubmitEvent) isEditorEvent() {}

type lineIterator struct {
	Lines     []text.Line
	Clip      image.Rectangle
	Alignment text.Alignment
	Width     int
	Offset    image.Point

	y, prevDesc fixed.Int26_6
}

func (l *lineIterator) Next() (text.String, f32.Point, bool) {
	for len(l.Lines) > 0 {
		line := l.Lines[0]
		l.Lines = l.Lines[1:]
		x := widgets.Align(l.Alignment, line.Width, l.Width) + fixed.I(l.Offset.X)
		l.y += l.prevDesc + line.Ascent
		l.prevDesc = line.Descent
		// Align baseline and line start to the pixel grid.
		off := fixed.Point26_6{X: fixed.I(x.Floor()), Y: fixed.I(l.y.Ceil())}
		l.y = off.Y
		off.Y += fixed.I(l.Offset.Y)
		if (off.Y + line.Bounds.Min.Y).Floor() > l.Clip.Max.Y {
			break
		}
		if (off.Y + line.Bounds.Max.Y).Ceil() < l.Clip.Min.Y {
			continue
		}
		str := line.Text
		for len(str.Advances) > 0 {
			adv := str.Advances[0]
			if (off.X + adv + line.Bounds.Max.X - line.Width).Ceil() >= l.Clip.Min.X {
				break
			}
			off.X += adv
			_, s := utf8.DecodeRuneInString(str.String)
			str.String = str.String[s:]
			str.Advances = str.Advances[1:]
		}
		n := 0
		endx := off.X
		for i, adv := range str.Advances {
			if (endx + line.Bounds.Min.X).Floor() > l.Clip.Max.X {
				str.String = str.String[:n]
				str.Advances = str.Advances[:i]
				break
			}
			_, s := utf8.DecodeRuneInString(str.String[n:])
			n += s
			endx += adv
		}
		offf := f32.Point{X: float32(off.X) / 64, Y: float32(off.Y) / 64}
		return str, offf, true
	}
	return text.String{}, f32.Point{}, false
}
