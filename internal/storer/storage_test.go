package storer

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStorer(t *testing.T) {
	t.Run("testGetEvents", func(t *testing.T) {
		s := New()
		firstGet, err := s.GetEvents()
		require.NotNil(t, firstGet)
		require.NoError(t, err)

		secondGet, err := s.GetEvents()
		require.NoError(t, err)
		require.NotNil(t, secondGet)
	})
}
