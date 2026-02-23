package elements

import "github.com/AnatoleLucet/go-opentui"

func (a *TextAreaElement) Placeholder() string {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.placeholder.Text()
}

func (a *TextAreaElement) SetPlaceholder(text string) {
	a.scheduleUpdate(func() error {
		a.mu.Lock()
		defer a.mu.Unlock()

		if err := guardDestroyed(a.ctx); err != nil {
			return err
		}

		a.placeholder.SetText(text)
		a.updatePlaceholder()
		return nil
	})
}

func (a *TextAreaElement) UnsetPlaceholder() {
	a.SetPlaceholder("")
}

func (a *TextAreaElement) SetPlaceholderFontWeight(weight string) {
	a.scheduleUpdate(func() error {
		a.mu.Lock()
		defer a.mu.Unlock()

		if err := guardDestroyed(a.ctx); err != nil {
			return err
		}

		a.placeholder.SetFontWeight(weight)
		a.updatePlaceholder()
		return nil
	})
}

// normal | bold
func (a *TextAreaElement) UnsetPlaceholderFontWeight() {
	a.scheduleUpdate(func() error {
		a.mu.Lock()
		defer a.mu.Unlock()

		if err := guardDestroyed(a.ctx); err != nil {
			return err
		}

		a.placeholder.UnsetFontWeight()
		a.updatePlaceholder()
		return nil
	})
}

// normal | italic
func (a *TextAreaElement) SetPlaceholderFontStyle(style string) {
	a.scheduleUpdate(func() error {
		a.mu.Lock()
		defer a.mu.Unlock()

		if err := guardDestroyed(a.ctx); err != nil {
			return err
		}

		a.placeholder.SetFontStyle(style)
		a.updatePlaceholder()
		return nil
	})
}

func (a *TextAreaElement) UnsetPlaceholderFontStyle() {
	a.scheduleUpdate(func() error {
		a.mu.Lock()
		defer a.mu.Unlock()

		if err := guardDestroyed(a.ctx); err != nil {
			return err
		}

		a.placeholder.UnsetFontStyle()
		a.updatePlaceholder()
		return nil
	})
}

// none | underline | strikethrough
func (a *TextAreaElement) SetPlaceholderDecoration(decoration string) {
	a.scheduleUpdate(func() error {
		a.mu.Lock()
		defer a.mu.Unlock()

		if err := guardDestroyed(a.ctx); err != nil {
			return err
		}

		a.placeholder.SetTextDecoration(decoration)
		a.updatePlaceholder()
		return nil
	})
}

func (a *TextAreaElement) UnsetPlaceholderDecoration() {
	a.scheduleUpdate(func() error {
		a.mu.Lock()
		defer a.mu.Unlock()

		if err := guardDestroyed(a.ctx); err != nil {
			return err
		}

		a.placeholder.UnsetTextDecoration()
		a.updatePlaceholder()
		return nil
	})
}

func (a *TextAreaElement) SetPlaceholderForeground(color string) {
	a.scheduleUpdate(func() error {
		a.mu.Lock()
		defer a.mu.Unlock()

		if err := guardDestroyed(a.ctx); err != nil {
			return err
		}

		a.placeholder.SetTextForeground(color)
		a.updatePlaceholder()
		return nil
	})
}

func (a *TextAreaElement) UnsetPlaceholderForeground() {
	a.scheduleUpdate(func() error {
		a.mu.Lock()
		defer a.mu.Unlock()

		if err := guardDestroyed(a.ctx); err != nil {
			return err
		}

		a.placeholder.UnsetTextForeground()
		a.updatePlaceholder()
		return nil
	})
}

func (a *TextAreaElement) SetPlaceholderBackground(color string) {
	a.scheduleUpdate(func() error {
		a.mu.Lock()
		defer a.mu.Unlock()

		if err := guardDestroyed(a.ctx); err != nil {
			return err
		}

		a.placeholder.SetTextBackground(color)
		a.updatePlaceholder()
		return nil
	})
}

func (a *TextAreaElement) UnsetPlaceholderBackground() {
	a.scheduleUpdate(func() error {
		a.mu.Lock()
		defer a.mu.Unlock()

		if err := guardDestroyed(a.ctx); err != nil {
			return err
		}

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
