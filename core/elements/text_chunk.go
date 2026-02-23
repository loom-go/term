package elements

import (
	"github.com/AnatoleLucet/go-opentui"
)

type TextChunk struct {
	text string

	bold          *bool
	italic        *bool
	underline     *bool
	strikethrough *bool

	foregroundColor *Color
	backgroundColor *Color
}

func NewTextChunk() *TextChunk {
	return &TextChunk{}
}

func (c *TextChunk) Text() string {
	return c.text
}

func (c *TextChunk) SetText(text string) {
	c.text = text
}

func (c *TextChunk) UnsetText() {
	c.SetText("")
}

// normal | bold
func (c *TextChunk) SetFontWeight(weight string) {
	if weight == "normal" {
		b := false
		c.bold = &b
	}
	if weight == "bold" {
		b := true
		c.bold = &b
	}
}

func (c *TextChunk) UnsetFontWeight() {
	c.bold = nil
}

// normal | italic
func (c *TextChunk) SetFontStyle(style string) {
	if style == "normal" {
		i := false
		c.italic = &i
	}
	if style == "italic" {
		i := true
		c.italic = &i
	}
}

func (c *TextChunk) UnsetFontStyle() {
	c.italic = nil
}

// none | underline | line-through
func (c *TextChunk) SetTextDecoration(decoration string) {
	if decoration == "none" {
		u := false
		s := false
		c.underline = &u
		c.strikethrough = &s
	}
	if decoration == "underline" {
		u := true
		s := false
		c.underline = &u
		c.strikethrough = &s
	}
	if decoration == "line-through" || decoration == "strike" {
		u := false
		s := true
		c.underline = &u
		c.strikethrough = &s
	}
}

func (c *TextChunk) UnsetTextDecoration() {
	c.underline = nil
	c.strikethrough = nil
}

func (c *TextChunk) SetTextForeground(color string) {
	if c.foregroundColor == nil {
		c.foregroundColor = NewColor("")
	}

	c.foregroundColor.Set(color)
}

func (c *TextChunk) UnsetTextForeground() {
	c.foregroundColor = nil
}

func (c *TextChunk) SetTextBackground(color string) {
	if c.backgroundColor == nil {
		c.backgroundColor = NewColor("")
	}

	c.backgroundColor.Set(color)
}

func (c *TextChunk) UnsetTextBackground() {
	c.backgroundColor = nil
}

func (c *TextChunk) StyledChunk(parent *TextChunk) opentui.StyledChunk {
	var chunk opentui.StyledChunk
	chunk.Text = c.text

	var parentBold *bool
	var parentItalic *bool
	var parentUnderline *bool
	var parentStrikethrough *bool
	var parentForeground *Color
	var parentBackground *Color

	if parent != nil {
		parentBold = parent.bold
		parentItalic = parent.italic
		parentUnderline = parent.underline
		parentStrikethrough = parent.strikethrough
		parentForeground = parent.foregroundColor
		parentBackground = parent.backgroundColor
	}

	if coalesce(c.bold, parentBold) {
		chunk.Attributes |= opentui.AttrBold
	}
	if coalesce(c.italic, parentItalic) {
		chunk.Attributes |= opentui.AttrItalic
	}
	if coalesce(c.underline, parentUnderline) {
		chunk.Attributes |= opentui.AttrUnderline
	}
	if coalesce(c.strikethrough, parentStrikethrough) {
		chunk.Attributes |= opentui.AttrStrike
	}

	if c.foregroundColor != nil {
		rgba := c.foregroundColor.RGBA()
		chunk.Foreground = &rgba
	} else if parentForeground != nil {
		rgba := parentForeground.RGBA()
		chunk.Foreground = &rgba
	}

	if c.backgroundColor != nil {
		rgba := c.backgroundColor.RGBA()
		chunk.Background = &rgba
	} else if parentBackground != nil {
		rgba := parentBackground.RGBA()
		chunk.Background = &rgba
	}

	return chunk
}

func coalesce[T any](values ...*T) T {
	var zero T

	for _, v := range values {
		if v != nil {
			return *v
		}
	}

	return zero
}
