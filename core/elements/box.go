package elements

import (
	"fmt"
	"github.com/AnatoleLucet/loom-term/core/gfx"

	"github.com/AnatoleLucet/go-opentui"
)

var borderStyles = map[string]opentui.BorderChars{
	"single":  opentui.BorderCharsSingle,
	"double":  opentui.BorderCharsDouble,
	"rounded": opentui.BorderCharsRounded,
	"heavy":   opentui.BorderCharsHeavy,
}

var textAlignments = map[string]opentui.TextAlignment{
	"left":   opentui.AlignLeft,
	"center": opentui.AlignCenter,
	"right":  opentui.AlignRight,
}

type BoxElement struct {
	*BaseElement

	title          string
	titleAlignment opentui.TextAlignment

	borderColor *Color
	borderChars opentui.BorderChars
	borderSides opentui.BorderSides

	backgroundColor *Color
}

func NewBoxElement() (box *BoxElement, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Box: %w: %v", ErrFailedToInitializeElement, err)
		}
	}()

	base, err := NewBaseElement()
	if err != nil {
		return nil, err
	}

	box = &BoxElement{
		BaseElement: base,

		borderColor:     NewColor("black"),
		backgroundColor: NewColor(""),
	}
	base.self = box

	return box, nil
}

func (e *BoxElement) SetTitle(title string) {
	e.scheduleUpdate(func() error {
		e.mu.Lock()
		defer e.mu.Unlock()

		if err := guardDestroyed(e.ctx); err != nil {
			return err
		}

		e.title = title
		return nil
	})
}

func (e *BoxElement) UnsetTitle() {
	e.scheduleUpdate(func() error {
		e.mu.Lock()
		defer e.mu.Unlock()

		if err := guardDestroyed(e.ctx); err != nil {
			return err
		}

		e.title = ""
		return nil
	})
}

// left | center | right
func (e *BoxElement) SetTitleAlignment(alignment string) {
	e.scheduleUpdate(func() error {
		e.mu.Lock()
		defer e.mu.Unlock()

		if err := guardDestroyed(e.ctx); err != nil {
			return err
		}

		align, ok := textAlignments[alignment]
		if !ok {
			return fmt.Errorf("invalid title alignment: %s", alignment)
		}

		e.titleAlignment = align
		return nil
	})
}

func (e *BoxElement) UnsetTitleAlignment() {
	e.scheduleUpdate(func() error {
		e.mu.Lock()
		defer e.mu.Unlock()

		if err := guardDestroyed(e.ctx); err != nil {
			return err
		}

		e.titleAlignment = opentui.AlignLeft
		return nil
	})
}

// single | double | rounded | heavy
func (e *BoxElement) SetBorderAll(style string) {
	e.scheduleUpdate(func() error {
		e.mu.Lock()
		defer e.mu.Unlock()

		if err := guardDestroyed(e.ctx); err != nil {
			return err
		}

		chars, ok := borderStyles[style]
		if !ok {
			return fmt.Errorf("invalid border style: %s", style)
		}

		e.borderChars = chars
		e.borderSides = opentui.BorderSides{Top: true, Right: true, Bottom: true, Left: true}
		return nil
	})
}

func (e *BoxElement) UnsetBorderAll() {
	e.scheduleUpdate(func() error {
		e.mu.Lock()
		defer e.mu.Unlock()

		if err := guardDestroyed(e.ctx); err != nil {
			return err
		}

		e.borderChars = opentui.BorderChars{}
		e.borderSides = opentui.BorderSides{}
		return nil
	})
}

// single | double | rounded | heavy
func (e *BoxElement) SetBorderHorizontal(style string) {
	e.scheduleUpdate(func() error {
		e.mu.Lock()
		defer e.mu.Unlock()

		if err := guardDestroyed(e.ctx); err != nil {
			return err
		}

		chars, ok := borderStyles[style]
		if !ok {
			return fmt.Errorf("invalid border style: %s", style)
		}

		e.borderChars = chars
		e.borderSides.Top = true
		e.borderSides.Bottom = true
		return nil
	})
}

// single | double | rounded | heavy
func (e *BoxElement) SetBorderVertical(style string) {
	e.scheduleUpdate(func() error {
		e.mu.Lock()
		defer e.mu.Unlock()

		if err := guardDestroyed(e.ctx); err != nil {
			return err
		}

		chars, ok := borderStyles[style]
		if !ok {
			return fmt.Errorf("invalid border style: %s", style)
		}

		e.borderChars = chars
		e.borderSides.Right = true
		e.borderSides.Left = true
		return nil
	})
}

// single | double | rounded | heavy
func (e *BoxElement) SetBorderTop(style string) {
	e.scheduleUpdate(func() error {
		e.mu.Lock()
		defer e.mu.Unlock()

		if err := guardDestroyed(e.ctx); err != nil {
			return err
		}

		chars, ok := borderStyles[style]
		if !ok {
			return fmt.Errorf("invalid border style: %s", style)
		}

		e.borderChars = chars
		e.borderSides.Top = true
		return nil
	})
}

// single | double | rounded | heavy
func (e *BoxElement) SetBorderRight(style string) {
	e.scheduleUpdate(func() error {
		e.mu.Lock()
		defer e.mu.Unlock()

		if err := guardDestroyed(e.ctx); err != nil {
			return err
		}

		chars, ok := borderStyles[style]
		if !ok {
			return fmt.Errorf("invalid border style: %s", style)
		}

		e.borderChars = chars
		e.borderSides.Right = true
		return nil
	})
}

// single | double | rounded | heavy
func (e *BoxElement) SetBorderBottom(style string) {
	e.scheduleUpdate(func() error {
		e.mu.Lock()
		defer e.mu.Unlock()

		if err := guardDestroyed(e.ctx); err != nil {
			return err
		}

		chars, ok := borderStyles[style]
		if !ok {
			return fmt.Errorf("invalid border style: %s", style)
		}

		e.borderChars = chars
		e.borderSides.Bottom = true
		return nil
	})
}

// single | double | rounded | heavy
func (e *BoxElement) SetBorderLeft(style string) {
	e.scheduleUpdate(func() error {
		e.mu.Lock()
		defer e.mu.Unlock()

		if err := guardDestroyed(e.ctx); err != nil {
			return err
		}

		chars, ok := borderStyles[style]
		if !ok {
			return fmt.Errorf("invalid border style: %s", style)
		}

		e.borderChars = chars
		e.borderSides.Left = true
		return nil
	})
}

// single | double | rounded | heavy
func (e *BoxElement) SetBorderColor(color string) {
	e.scheduleUpdate(func() error {
		e.mu.Lock()
		defer e.mu.Unlock()

		if err := guardDestroyed(e.ctx); err != nil {
			return err
		}

		return e.borderColor.Set(color)
	})
}

func (e *BoxElement) UnsetBorderColor() {
	e.scheduleUpdate(func() error {
		e.mu.Lock()
		defer e.mu.Unlock()

		if err := guardDestroyed(e.ctx); err != nil {
			return err
		}

		return e.borderColor.Set("transparent")
	})
}

// single | double | rounded | heavy
func (e *BoxElement) SetBackgroundColor(color string) {
	e.scheduleUpdate(func() error {
		e.mu.Lock()
		defer e.mu.Unlock()

		if err := guardDestroyed(e.ctx); err != nil {
			return err
		}

		return e.backgroundColor.Set(color)
	})
}

func (e *BoxElement) UnsetBackgroundColor() {
	e.scheduleUpdate(func() error {
		e.mu.Lock()
		defer e.mu.Unlock()

		if err := guardDestroyed(e.ctx); err != nil {
			return err
		}

		return e.backgroundColor.Set("transparent")
	})
}

func (b *BoxElement) Render(buffer *opentui.Buffer, rect gfx.Rect) error {
	borderColor := b.borderColor.RGBA()
	backgroundColor := b.backgroundColor.RGBA()

	opts := opentui.BoxOptions{
		Fill: backgroundColor.A > 0,

		Title:          b.title,
		TitleAlignment: b.titleAlignment,

		Sides:       b.borderSides,
		BorderChars: b.borderChars,
	}

	err := buffer.DrawBox(
		int32(rect.X), int32(rect.Y),
		uint32(rect.Width), uint32(rect.Height),
		opts,
		borderColor,
		backgroundColor,
	)

	if err != nil {
		return fmt.Errorf("Box: %w", err)
	}

	return nil
}
