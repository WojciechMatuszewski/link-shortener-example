package link_test

import (
	"testing"

	"link-shortener/src/internal/platform/link"

	"github.com/stretchr/testify/require"
)

func TestIsOriginValid(t *testing.T) {
	t.Run("malformed", func(t *testing.T) {
		require.False(t, link.IsOriginValid("foo"))
	})
	t.Run("valid", func(t *testing.T) {
		require.True(t, link.IsOriginValid("https://google.com/something?q=123"))
	})
}

func TestNew(t *testing.T) {
	t.Run("correctly creates an link", func(t *testing.T) {
		service := link.NewServie("https://wojtek.tk")
		l, err := service.New("https://google.com")
		require.NoError(t, err)
		require.Contains(t, l.Path, "https://wojtek.tk")
	})

	t.Run("validates origin", func(t *testing.T) {
		service := link.NewServie("https://wojtek.tk")
		_, err := service.New("foo")
		require.Equal(t, link.ErrOriginNotValid, err)
	})
}
