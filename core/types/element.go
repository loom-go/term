package types

import (
	"iter"

	"github.com/AnatoleLucet/go-opentui"
	"github.com/AnatoleLucet/tess"
)

type Element interface {
	Parent() Element
	SetParent(parent Element) error

	Children() iter.Seq[Element]
	AppendChild(child Element) error
	RemoveChild(child Element) error

	HandleMouseMove(*EventMouse) error
	HandleMouseEnter(*EventMouse) error
	HandleMouseLeave(*EventMouse) error
	HandleMousePress(*EventMouse) error
	HandleMouseRelease(*EventMouse) error
	HandleMouseScroll(*EventMouse) error
	HandleMouseDrag(*EventMouse) error

	OnMouseMove(handler func(*EventMouse)) (remove func())
	OnMouseEnter(handler func(*EventMouse)) (remove func())
	OnMouseLeave(handler func(*EventMouse)) (remove func())
	OnMousePress(handler func(*EventMouse)) (remove func())
	OnMouseRelease(handler func(*EventMouse)) (remove func())
	OnMouseScroll(handler func(*EventMouse)) (remove func())
	OnMouseDrag(handler func(*EventMouse)) (remove func())

	Paint(buffer *opentui.Buffer, x, y float32) error

	XYZ() *tess.Node
	Layout() *tess.Layout

	Destroy() error

	ZIndex() int
	SetZIndex(zIndex int) error
	UnsetZIndex() error

	// 10 | "100pt" | "50%" | "auto" | "max-content" | "fit-content" | "stretch"
	SetWidth(width any) error
	UnsetWidth() error
	// 10 | "100pt" | "50%" | "auto" | "max-content" | "fit-content" | "stretch"
	SetMinWidth(minWidth any) error
	UnsetMinWidth() error
	// 10 | "100pt" | "50%" | "auto" | "max-content" | "fit-content" | "stretch"
	SetMaxWidth(maxWidth any) error
	UnsetMaxWidth() error
	// 10 | "100pt" | "50%" | "auto" | "max-content" | "fit-content" | "stretch"
	SetHeight(height any) error
	UnsetHeight() error
	// 10 | "100pt" | "50%" | "auto" | "max-content" | "fit-content" | "stretch"
	SetMinHeight(minHeight any) error
	UnsetMinHeight() error
	// 10 | "100pt" | "50%" | "auto" | "max-content" | "fit-content" | "stretch"
	SetMaxHeight(maxHeight any) error
	UnsetMaxHeight() error

	// 10 | "100pt" | "50%" | "auto" | "max-content" | "fit-content" | "stretch"
	SetTop(top any) error
	UnsetTop() error
	// 10 | "100pt" | "50%" | "auto" | "max-content" | "fit-content" | "stretch"
	SetBottom(bottom any) error
	UnsetBottom() error
	// 10 | "100pt" | "50%" | "auto" | "max-content" | "fit-content" | "stretch"
	SetLeft(left any) error
	UnsetLeft() error
	// 10 | "100pt" | "50%" | "auto" | "max-content" | "fit-content" | "stretch"
	SetRight(right any) error
	UnsetRight() error

	// "static" | "relative" | "absolute"
	SetPosition(position string) error
	UnsetPosition() error

	// 10 | "100pt" | "50%"
	SetPaddingAll(padding any) error
	UnsetPaddingAll() error
	// 10 | "100pt" | "50%"
	SetPaddingVertical(padding any) error
	UnsetPaddingVertical() error
	// 10 | "100pt" | "50%"
	SetPaddingHorizontal(padding any) error
	UnsetPaddingHorizontal() error
	// 10 | "100pt" | "50%"
	SetPaddingTop(padding any) error
	UnsetPaddingTop() error
	// 10 | "100pt" | "50%"
	SetPaddingBottom(padding any) error
	UnsetPaddingBottom() error
	// 10 | "100pt" | "50%"
	SetPaddingLeft(padding any) error
	UnsetPaddingLeft() error
	// 10 | "100pt" | "50%"
	SetPaddingRight(padding any) error
	UnsetPaddingRight() error

	// 10 | "100pt" | "50%"
	SetMarginAll(margin any) error
	UnsetMarginAll() error
	// 10 | "100pt" | "50%"
	SetMarginVertical(margin any) error
	UnsetMarginVertical() error
	// 10 | "100pt" | "50%"
	SetMarginHorizontal(margin any) error
	UnsetMarginHorizontal() error
	// 10 | "100pt" | "50%"
	SetMarginTop(margin any) error
	UnsetMarginTop() error
	// 10 | "100pt" | "50%"
	SetMarginBottom(margin any) error
	UnsetMarginBottom() error
	// 10 | "100pt" | "50%"
	SetMarginLeft(margin any) error
	UnsetMarginLeft() error
	// 10 | "100pt" | "50%"
	SetMarginRight(margin any) error
	UnsetMarginRight() error

	// "none" | "flex" | "contents"
	SetDisplay(display string) error
	UnsetDisplay() error

	// "start" | "end" | "center" | "stretch" | "baseline"
	SetAlignSelf(alignSelf string) error
	UnsetAlignSelf() error
	// "start" | "end" | "center
	SetAlignItems(alignItems string) error
	UnsetAlignItems() error
	// "start" | "end" | "center" | "stretch" | "baseline"
	SetAlignContent(alignContent string) error
	UnsetAlignContent() error
	// "start" | "end" | "center" | "space-between" | "space-around" | "space-evenly"
	SetJustifyContent(justifyContent string) error
	UnsetJustifyContent() error
	// "row" | "row-reverse" | "column" | "column-reverse"
	SetFlexDirection(flexDirection string) error
	UnsetFlexDirection() error
	// "nowrap" | "wrap" | "wrap-reverse"
	SetFlexWrap(flexWrap string) error
	UnsetFlexWrap() error
	// "none" | "0" | "1" | ...
	SetFlexGrow(flexGrow string) error
	UnsetFlexGrow() error
	// "none" | "0" | "1" | ...
	SetFlexShrink(flexShrink string) error
	UnsetFlexShrink() error

	// 10 | "100pt" | "50%"
	SetGapAll(gap any) error
	UnsetGapAll() error
	// 10 | "100pt" | "50%"
	SetGapRow(gap any) error
	UnsetGapRow() error
	// 10 | "100pt" | "50%"
	SetGapColumn(gap any) error
	UnsetGapColumn() error

	// "visible" | "hidden"
	SetOverflow(overflow string) error
	UnsetOverflow() error

	// todo: doens't belong here. should be on box element itself
	// // "single" | "double" | "rounded"  | "heavy"
	// SetBorderStyle(borderStyle string) error
}
