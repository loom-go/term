package elements

import "github.com/AnatoleLucet/go-opentui"

func (a *TextAreaElement) Placeholder() (text string) {
	scheduleAccess(a.Self(), func() {
		a.mu.RLock()
		defer a.mu.RUnlock()

		text = a.placeholder.Text()
	})

	return
}

func (a *TextAreaElement) SetPlaceholder(text string) {
	scheduleUpdate(a.Self(), func() error {
		a.mu.Lock()
		defer a.mu.Unlock()

		a.placeholder.SetText(text)
		a.updatePlaceholder()
		return nil
	})
}

func (a *TextAreaElement) UnsetPlaceholder() {
	a.SetPlaceholder("")
}

func (a *TextAreaElement) SetPlaceholderFontWeight(weight string) {
	scheduleUpdate(a.Self(), func() error {
		a.mu.Lock()
		defer a.mu.Unlock()

		a.placeholder.SetFontWeight(weight)
		a.updatePlaceholder()
		return nil
	})
}

// normal | bold
func (a *TextAreaElement) UnsetPlaceholderFontWeight() {
	scheduleUpdate(a.Self(), func() error {
		a.mu.Lock()
		defer a.mu.Unlock()

		a.placeholder.UnsetFontWeight()
		a.updatePlaceholder()
		return nil
	})
}

// normal | italic
func (a *TextAreaElement) SetPlaceholderFontStyle(style string) {
	scheduleUpdate(a.Self(), func() error {
		a.mu.Lock()
		defer a.mu.Unlock()

		a.placeholder.SetFontStyle(style)
		a.updatePlaceholder()
		return nil
	})
}

func (a *TextAreaElement) UnsetPlaceholderFontStyle() {
	scheduleUpdate(a.Self(), func() error {
		a.mu.Lock()
		defer a.mu.Unlock()

		a.placeholder.UnsetFontStyle()
		a.updatePlaceholder()
		return nil
	})
}

// none | underline | strikethrough
func (a *TextAreaElement) SetPlaceholderDecoration(decoration string) {
	scheduleUpdate(a.Self(), func() error {
		a.mu.Lock()
		defer a.mu.Unlock()

		a.placeholder.SetTextDecoration(decoration)
		a.updatePlaceholder()
		return nil
	})
}

func (a *TextAreaElement) UnsetPlaceholderDecoration() {
	scheduleUpdate(a.Self(), func() error {
		a.mu.Lock()
		defer a.mu.Unlock()

		a.placeholder.UnsetTextDecoration()
		a.updatePlaceholder()
		return nil
	})
}

func (a *TextAreaElement) SetPlaceholderForeground(color string) {
	scheduleUpdate(a.Self(), func() error {
		a.mu.Lock()
		defer a.mu.Unlock()

		a.placeholder.SetTextForeground(color)
		a.updatePlaceholder()
		return nil
	})
}

func (a *TextAreaElement) UnsetPlaceholderForeground() {
	scheduleUpdate(a.Self(), func() error {
		a.mu.Lock()
		defer a.mu.Unlock()

		a.placeholder.UnsetTextForeground()
		a.updatePlaceholder()
		return nil
	})
}

func (a *TextAreaElement) SetPlaceholderBackground(color string) {
	scheduleUpdate(a.Self(), func() error {
		a.mu.Lock()
		defer a.mu.Unlock()

		a.placeholder.SetTextBackground(color)
		a.updatePlaceholder()
		return nil
	})
}

func (a *TextAreaElement) UnsetPlaceholderBackground() {
	scheduleUpdate(a.Self(), func() error {
		a.mu.Lock()
		defer a.mu.Unlock()

		a.placeholder.UnsetTextBackground()
		a.updatePlaceholder()
		return nil
	})
}

func (a *TextAreaElement) updatePlaceholder() {
	var chunks []opentui.StyledChunk
	if a.placeholder.Text() != "" {
		chunks = append(chunks, a.placeholder.StyledChunk(nil))
	}

	a.editBufferView.SetPlaceholderStyledText(chunks)
}
