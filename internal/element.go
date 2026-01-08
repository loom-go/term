package internal

import (
	"slices"
	"sync"

	"github.com/AnatoleLucet/loom-term/opentui"
	"github.com/AnatoleLucet/tess"
)

type Element struct {
	mu sync.RWMutex

	layout *tess.Node

	parent   *Element
	children []*Element

	paint func(self *Element, buffer *opentui.Buffer) error

	// Render data - these get cloned for snapshot rendering
	backgroundColor opentui.RGBA
	text            string
	textColor       opentui.RGBA
}

func NewElement(parent *Element) (*Element, error) {
	layout, err := tess.NewNode()
	if err != nil {
		return nil, err
	}

	return &Element{
		layout: layout,
		parent: parent,
	}, nil
}

// todo: find a better name as it's conflicting with tess' layout naming
func (n *Element) Layout() *tess.Node {
	return n.layout
}

func (n *Element) Parent() *Element {
	return n.parent
}

func (n *Element) Children() []*Element {
	return n.children
}

func (n *Element) AppendChild(child *Element) error {
	n.children = append(n.children, child)
	n.layout.AddChild(child.layout)

	return nil
}

func (n *Element) RemoveChild(child *Element) error {
	i := slices.Index(n.children, child)
	n.children = append(n.children[:i], n.children[i+1:]...)
	n.layout.RemoveChild(child.layout)

	return nil
}

func (n *Element) Paint(buffer *opentui.Buffer) error {
	n.mu.RLock()
	paint := n.paint
	n.mu.RUnlock()

	if paint != nil {
		return paint(n, buffer)
	}

	return nil
}

func (n *Element) PaintAll(buffer *opentui.Buffer) error {
	n.Paint(buffer)

	for _, child := range n.children {
		child.PaintAll(buffer)
	}

	return nil
}

func (n *Element) SetPaint(paintFn func(self *Element, buffer *opentui.Buffer) error) {
	n.mu.Lock()
	n.paint = paintFn
	n.mu.Unlock()
}

// Setters for render data (with locking)

func (n *Element) SetBackgroundColor(color opentui.RGBA) {
	n.mu.Lock()
	n.backgroundColor = color
	n.mu.Unlock()
}

func (n *Element) SetText(text string) {
	n.mu.Lock()
	n.text = text
	n.mu.Unlock()
}

func (n *Element) SetTextColor(color opentui.RGBA) {
	n.mu.Lock()
	n.textColor = color
	n.mu.Unlock()
}

// Getters for render data (with locking) - used by Paint functions on cloned elements

func (n *Element) BackgroundColor() opentui.RGBA {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.backgroundColor
}

func (n *Element) Text() string {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.text
}

func (n *Element) TextColor() opentui.RGBA {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.textColor
}

func (n *Element) Clone() *Element {
	n.mu.RLock()
	clone := &Element{
		paint:           n.paint,
		backgroundColor: n.backgroundColor,
		text:            n.text,
		textColor:       n.textColor,
	}
	n.mu.RUnlock()

	clone.layout = n.layout.Clone()
	clone.layout.RemoveAllChildren()

	clone.children = make([]*Element, len(n.children))
	for i, child := range n.children {
		childClone := child.Clone()
		childClone.parent = clone

		clone.children[i] = childClone
		clone.layout.AddChild(childClone.layout)
	}

	return clone
}
