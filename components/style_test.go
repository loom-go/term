package components

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRGB(t *testing.T) {
	t.Run("valid RGBA hex", func(t *testing.T) {
		color := RGBA(100, 150, 200, 0.75)

		assert.Equal(t, "#6496c8bf", color)
	})

	t.Run("valid RGB hex", func(t *testing.T) {
		color := RGBA(255, 0, 128, 1)

		assert.Equal(t, "#ff0080ff", color)
	})
}
