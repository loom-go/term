package elements

import (
	"github.com/AnatoleLucet/loom-term/core/gfx"
	"iter"

	"github.com/AnatoleLucet/tess"
)

type Element interface {
	ElementTree
	ElementStyle
	ElementEvent

	Destroy()
	destroyUnsafe()

	gfx.Renderable
	gfx.Recordable
}

type ElementTree interface {
	lock()
	unlock()
	setContextUnsafe(*RenderContext)
	flushPendingUpdates() error

	ID() uint32

	Self() Element

	Parent() Element
	parentUnsafe() Element
	SetParent(Element)
	setParentUnsafe(Element) error

	Children() iter.Seq[Element]
	childrenUnsafe() iter.Seq[Element]
	AppendChild(Element)
	appendChildUnsafe(Element) error
	RemoveChild(Element)
	removeChildUnsafe(Element) error
}

type ElementStyle interface {
	xyz() *tess.Node

	ZIndex() int
	zindexUnsafe() int

	SetZIndex(zIndex int)
	UnsetZIndex()
	updateZIndexUnsafe(child Element, oldz, newz int) error

	// 10 | "100pt" | "50%" | "auto" | "max-content" | "fit-content" | "stretch"
	SetWidth(width any)
	UnsetWidth()
	// 10 | "100pt" | "50%" | "auto" | "max-content" | "fit-content" | "stretch"
	SetMinWidth(minWidth any)
	UnsetMinWidth()
	// 10 | "100pt" | "50%" | "auto" | "max-content" | "fit-content" | "stretch"
	SetMaxWidth(maxWidth any)
	UnsetMaxWidth()
	// 10 | "100pt" | "50%" | "auto" | "max-content" | "fit-content" | "stretch"
	SetHeight(height any)
	UnsetHeight()
	// 10 | "100pt" | "50%" | "auto" | "max-content" | "fit-content" | "stretch"
	SetMinHeight(minHeight any)
	UnsetMinHeight()
	// 10 | "100pt" | "50%" | "auto" | "max-content" | "fit-content" | "stretch"
	SetMaxHeight(maxHeight any)
	UnsetMaxHeight()

	SetTranslate(x, y float32)
	UnsetTranslate()
	// 10 | "100pt" | "50%" | "auto" | "max-content" | "fit-content" | "stretch"
	SetTop(top any)
	UnsetTop()
	// 10 | "100pt" | "50%" | "auto" | "max-content" | "fit-content" | "stretch"
	SetBottom(bottom any)
	UnsetBottom()
	// 10 | "100pt" | "50%" | "auto" | "max-content" | "fit-content" | "stretch"
	SetLeft(left any)
	UnsetLeft()
	// 10 | "100pt" | "50%" | "auto" | "max-content" | "fit-content" | "stretch"
	SetRight(right any)
	UnsetRight()

	// "static" | "relative" | "absolute"
	SetPosition(position string)
	UnsetPosition()

	// 10 | "100pt" | "50%"
	SetPaddingAll(padding any)
	UnsetPaddingAll()
	// 10 | "100pt" | "50%"
	SetPaddingVertical(padding any)
	UnsetPaddingVertical()
	// 10 | "100pt" | "50%"
	SetPaddingHorizontal(padding any)
	UnsetPaddingHorizontal()
	// 10 | "100pt" | "50%"
	SetPaddingTop(padding any)
	UnsetPaddingTop()
	// 10 | "100pt" | "50%"
	SetPaddingBottom(padding any)
	UnsetPaddingBottom()
	// 10 | "100pt" | "50%"
	SetPaddingLeft(padding any)
	UnsetPaddingLeft()
	// 10 | "100pt" | "50%"
	SetPaddingRight(padding any)
	UnsetPaddingRight()

	// 10 | "100pt" | "50%"
	SetMarginAll(margin any)
	UnsetMarginAll()
	// 10 | "100pt" | "50%"
	SetMarginVertical(margin any)
	UnsetMarginVertical()
	// 10 | "100pt" | "50%"
	SetMarginHorizontal(margin any)
	UnsetMarginHorizontal()
	// 10 | "100pt" | "50%"
	SetMarginTop(margin any)
	UnsetMarginTop()
	// 10 | "100pt" | "50%"
	SetMarginBottom(margin any)
	UnsetMarginBottom()
	// 10 | "100pt" | "50%"
	SetMarginLeft(margin any)
	UnsetMarginLeft()
	// 10 | "100pt" | "50%"
	SetMarginRight(margin any)
	UnsetMarginRight()

	// "none" | "flex" | "contents"
	SetDisplay(display string)
	UnsetDisplay()

	// "start" | "end" | "center" | "stretch" | "baseline"
	SetAlignSelf(alignSelf string)
	UnsetAlignSelf()
	// "start" | "end" | "center
	SetAlignItems(alignItems string)
	UnsetAlignItems()
	// "start" | "end" | "center" | "stretch" | "baseline"
	SetAlignContent(alignContent string)
	UnsetAlignContent()
	// "start" | "end" | "center" | "space-between" | "space-around" | "space-evenly"
	SetJustifyContent(justifyContent string)
	UnsetJustifyContent()
	// "row" | "row-reverse" | "column" | "column-reverse"
	SetFlexDirection(flexDirection string)
	UnsetFlexDirection()
	// "nowrap" | "wrap" | "wrap-reverse"
	SetFlexWrap(flexWrap string)
	UnsetFlexWrap()
	// "none" | "0" | "1" | ...
	SetFlexGrow(flexGrow string)
	UnsetFlexGrow()
	// "none" | "0" | "1" | ...
	SetFlexShrink(flexShrink string)
	UnsetFlexShrink()

	// 10 | "100pt" | "50%"
	SetGapAll(gap any)
	UnsetGapAll()
	// 10 | "100pt" | "50%"
	SetGapRow(gap any)
	UnsetGapRow()
	// 10 | "100pt" | "50%"
	SetGapColumn(gap any)
	UnsetGapColumn()

	// "visible" | "hidden"
	SetOverflow(overflow string)
	UnsetOverflow()
}

type ElementEvent interface {
	setFocused(bool)
	broadcastEvent(EventType, any)

	OnMousePress(func(*EventMouse), ...EventOptions) (remove func())
	OnMouseRelease(func(*EventMouse), ...EventOptions) (remove func())
	OnMouseMove(func(*EventMouse), ...EventOptions) (remove func())
	OnMouseScroll(func(*EventMouse), ...EventOptions) (remove func())
	OnMouseDrag(func(*EventMouse), ...EventOptions) (remove func())
	OnMouseEnter(func(*EventMouse)) (remove func())
	OnMouseLeave(func(*EventMouse)) (remove func())

	OnKeyPress(func(*EventKey), ...EventOptions) (remove func())
	OnKeyRelease(func(*EventKey), ...EventOptions) (remove func())

	OnPaste(func(*EventPaste), ...EventOptions) (remove func())

	Focus()
	OnFocus(func(*EventFocus), ...EventOptions) (remove func())

	Blur()
	OnBlur(func(*EventBlur), ...EventOptions) (remove func())

	OnSubmit(func(*EventSubmit), ...EventOptions) (remove func())

	OnDestroy(func()) (remove func())
}
