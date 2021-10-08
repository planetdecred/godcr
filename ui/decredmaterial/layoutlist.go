package decredmaterial

import (
	"gioui.org/layout"
	"image"

	"gioui.org/gesture"
	"gioui.org/io/pointer"
	"gioui.org/op"
	"gioui.org/op/clip"
)

type scrollChild struct {
	size image.Point
	call op.CallOp
}

// List displays a subsection of a potentially infinitely
// large underlying list. List accepts user input to scroll
// the subsection.
type LayoutList struct {
	Axis layout.Axis
	// ScrollToEnd instructs the list to stay scrolled to the far end position
	// once reached. A List with ScrollToEnd == true and LayoutPosition.BeforeEnd ==
	// false draws its content with the last item at the bottom of the list
	// area.
	ScrollToEnd bool
	// Alignment is the cross axis alignment of list elements.
	Alignment layout.Alignment

	cs          layout.Constraints
	scroll      gesture.Scroll
	scrollDelta int

	// LayoutPosition is updated during Layout. To save the list scroll position,
	// just save LayoutPosition after Layout finishes. To scroll the list
	// programmatically, update LayoutPosition (e.g. restore it from a saved value)
	// before calling Layout.
	LayoutPosition LayoutPosition

	len int

	// maxSize is the total size of visible children.
	maxSize  int
	children []scrollChild
	dir      iterationDir
}

// ListElement is a function that computes the layout.Dimensions of
// a list element.
type ListElement func(gtx layout.Context, index int) layout.Dimensions

type iterationDir uint8

// LayoutPosition is a List scroll offset represented as an offset from the top edge
// of a child element.
type LayoutPosition struct {
	// BeforeEnd tracks whether the List position is before the very end. We
	// use "before end" instead of "at end" so that the zero value of a
	// LayoutPosition struct is useful.
	//
	// When laying out a list, if ScrollToEnd is true and BeforeEnd is false,
	// then First and Offset are ignored, and the list is drawn with the last
	// item at the bottom. If ScrollToEnd is false then BeforeEnd is ignored.
	BeforeEnd bool
	// First is the index of the first visible child.
	First int
	// Offset is the distance in pixels from the top edge to the child at index
	// First.
	Offset int
	// OffsetLast is the signed distance in pixels from the bottom edge to the
	// bottom edge of the child at index First+Count.
	OffsetLast int
	// Count is the number of visible children.
	Count int
	// Length is the estimated total size of all children, measured in pixels.
	Length int
}

const (
	iterateNone iterationDir = iota
	iterateForward
	iterateBackward
)

// init prepares the list for iterating through its children with next.
func (l *LayoutList) init(gtx layout.Context, len int) {
	if l.more() {
		panic("unfinished child")
	}
	l.cs = gtx.Constraints
	l.maxSize = 0
	l.children = l.children[:0]
	l.len = len
	l.update(gtx)
	if l.scrollToEnd() || l.LayoutPosition.First > len {
		l.LayoutPosition.Offset = 0
		l.LayoutPosition.First = len
	}
}

func crossConstraint(ax layout.Axis, cs layout.Constraints) (int, int) {
	if ax == layout.Horizontal {
		return cs.Min.Y, cs.Max.Y
	}
	return cs.Min.X, cs.Max.X
}

// constraints returns the constraints for axis a.
func constraints(ax layout.Axis, mainMin, mainMax, crossMin, crossMax int) layout.Constraints {
	if ax == layout.Horizontal {
		return layout.Constraints{Min: image.Pt(mainMin, crossMin), Max: image.Pt(mainMax, crossMax)}
	}
	return layout.Constraints{Min: image.Pt(crossMin, mainMin), Max: image.Pt(crossMax, mainMax)}
}

// Layout the List.
func (l *LayoutList) Layout(gtx layout.Context, len int, w ListElement) layout.Dimensions {
	l.init(gtx, len)
	crossMin, crossMax := crossConstraint(l.Axis, gtx.Constraints)
	gtx.Constraints = constraints(l.Axis, 0, inf, crossMin, crossMax)
	macro := op.Record(gtx.Ops)
	laidOutTotalLength := 0
	numLaidOut := 0

	for l.next(); l.more(); l.next() {
		child := op.Record(gtx.Ops)
		dims := w(gtx, l.index())
		call := child.Stop()
		l.end(dims, call)
		laidOutTotalLength += l.Axis.Convert(dims.Size).X
		numLaidOut++
	}

	if numLaidOut > 0 {
		l.LayoutPosition.Length = laidOutTotalLength * len / numLaidOut
	} else {
		l.LayoutPosition.Length = 0
	}
	return l.layout(gtx.Ops, macro)
}

func (l *LayoutList) scrollToEnd() bool {
	return l.ScrollToEnd && !l.LayoutPosition.BeforeEnd
}

// Dragging reports whether the List is being dragged.
func (l *LayoutList) Dragging() bool {
	return l.scroll.State() == gesture.StateDragging
}

func (l *LayoutList) update(gtx layout.Context) {
	d := l.scroll.Scroll(gtx.Metric, gtx, gtx.Now, gesture.Axis(l.Axis))
	l.scrollDelta = d
	l.LayoutPosition.Offset += d
}

// next advances to the next child.
func (l *LayoutList) next() {
	l.dir = l.nextDir()
	// The user scroll offset is applied after scrolling to
	// list end.
	if l.scrollToEnd() && !l.more() && l.scrollDelta < 0 {
		l.LayoutPosition.BeforeEnd = true
		l.LayoutPosition.Offset += l.scrollDelta
		l.dir = l.nextDir()
	}
}

// index is current child's position in the underlying list.
func (l *LayoutList) index() int {
	switch l.dir {
	case iterateBackward:
		return l.LayoutPosition.First - 1
	case iterateForward:
		return l.LayoutPosition.First + len(l.children)
	default:
		panic("Index called before Next")
	}
}

// more reports whether more children are needed.
func (l *LayoutList) more() bool {
	return l.dir != iterateNone
}

func mainConstraint(ax layout.Axis, cs layout.Constraints) (int, int) {
	if ax == layout.Horizontal {
		return cs.Min.X, cs.Max.X
	}
	return cs.Min.Y, cs.Max.Y
}

func (l *LayoutList) nextDir() iterationDir {

	_, vsize := mainConstraint(l.Axis, l.cs)
	last := l.LayoutPosition.First + len(l.children)
	// Clamp offset.
	if l.maxSize-l.LayoutPosition.Offset < vsize && last == l.len {
		l.LayoutPosition.Offset = l.maxSize - vsize
	}
	if l.LayoutPosition.Offset < 0 && l.LayoutPosition.First == 0 {
		l.LayoutPosition.Offset = 0
	}
	switch {
	case len(l.children) == l.len:
		return iterateNone
	case l.maxSize-l.LayoutPosition.Offset < vsize:
		return iterateForward
	case l.LayoutPosition.Offset < 0:
		return iterateBackward
	}
	return iterateNone
}

// End the current child by specifying its layout.Dimensions.
func (l *LayoutList) end(dims layout.Dimensions, call op.CallOp) {
	child := scrollChild{dims.Size, call}
	mainSize := l.Axis.Convert(child.size).X
	l.maxSize += mainSize
	switch l.dir {
	case iterateForward:
		l.children = append(l.children, child)
	case iterateBackward:
		l.children = append(l.children, scrollChild{})
		copy(l.children[1:], l.children)
		l.children[0] = child
		l.LayoutPosition.First--
		l.LayoutPosition.Offset += mainSize
	default:
		panic("call Next before End")
	}
	l.dir = iterateNone
}

// Layout the List and return its layout.Dimensions.
func (l *LayoutList) layout(ops *op.Ops, macro op.MacroOp) layout.Dimensions {
	if l.more() {
		panic("unfinished child")
	}
	mainMin, mainMax := mainConstraint(l.Axis, l.cs)
	children := l.children
	// Skip invisible children
	for len(children) > 0 {
		sz := children[0].size
		mainSize := l.Axis.Convert(sz).X
		if l.LayoutPosition.Offset < mainSize {
			// First child is partially visible.
			break
		}
		l.LayoutPosition.First++
		l.LayoutPosition.Offset -= mainSize
		children = children[1:]
	}
	size := -l.LayoutPosition.Offset
	var maxCross int
	for i, child := range children {
		sz := l.Axis.Convert(child.size)
		if c := sz.Y; c > maxCross {
			maxCross = c
		}
		size += sz.X
		if size >= mainMax {
			children = children[:i+1]
			break
		}
	}
	l.LayoutPosition.Count = len(children)
	l.LayoutPosition.OffsetLast = mainMax - size
	pos := -l.LayoutPosition.Offset
	// ScrollToEnd lists are end aligned.
	if space := l.LayoutPosition.OffsetLast; l.ScrollToEnd && space > 0 {
		pos += space
	}
	for _, child := range children {
		sz := l.Axis.Convert(child.size)
		var cross int
		switch l.Alignment {
		case layout.End:
			cross = maxCross - sz.Y
		case layout.Middle:
			cross = (maxCross - sz.Y) / 2
		}
		childSize := sz.X
		max := childSize + pos
		if max > mainMax {
			max = mainMax
		}
		min := pos
		if min < 0 {
			min = 0
		}
		r := image.Rectangle{
			Min: l.Axis.Convert(image.Pt(min, -inf)),
			Max: l.Axis.Convert(image.Pt(max, inf)),
		}
		stack := op.Save(ops)
		clip.Rect(r).Add(ops)
		pt := l.Axis.Convert(image.Pt(pos, cross))
		op.Offset(layout.FPt(pt)).Add(ops)
		child.call.Add(ops)
		stack.Load()
		pos += childSize
	}
	atStart := l.LayoutPosition.First == 0 && l.LayoutPosition.Offset <= 0
	atEnd := l.LayoutPosition.First+len(children) == l.len && mainMax >= pos
	if atStart && l.scrollDelta < 0 || atEnd && l.scrollDelta > 0 {
		l.scroll.Stop()
	}
	l.LayoutPosition.BeforeEnd = !atEnd
	if pos < mainMin {
		pos = mainMin
	}
	if pos > mainMax {
		pos = mainMax
	}
	dims := l.Axis.Convert(image.Pt(pos, maxCross))
	call := macro.Stop()
	defer op.Save(ops).Load()
	pointer.Rect(image.Rectangle{Max: dims}).Add(ops)

	var min, max int
	if o := l.LayoutPosition.Offset; o > 0 {
		// Use the size of the invisible part as scroll boundary.
		min = -o
	} else if l.LayoutPosition.First > 0 {
		min = -inf
	}
	if o := l.LayoutPosition.OffsetLast; o < 0 {
		max = -o
	} else if l.LayoutPosition.First+l.LayoutPosition.Count < l.len {
		max = inf
	}
	scrollRange := image.Rectangle{
		Min: l.Axis.Convert(image.Pt(min, 0)),
		Max: l.Axis.Convert(image.Pt(max, 0)),
	}
	l.scroll.Add(ops, scrollRange)

	call.Add(ops)
	return layout.Dimensions{Size: dims}
}
