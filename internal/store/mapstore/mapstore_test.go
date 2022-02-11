package mapstore

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSetGet(t *testing.T) {
	m := NewLinkMap(10)

	cases := map[string]struct {
		key      uint32
		shortURL string
		longURL  string
	}{
		"example": {
			4255220111,
			"x6Qhseaaaa",
			"https://www.youtube.com/watch?v=dQw4w9WgXcQ",
		},
		"empty one": {
			0,
			"",
			"",
		},
		"getError": {
			3,
			"",
			"",
		},
		"setError": {
			4255220111,
			"x6Qhseaaaa",
			"https://www.youtube.com/watch?v=dQw4w9WgXcQ",
		},
	}
	var temp string
	for name, c := range cases {
		if name == "getError" {
			temp, err := m.Get(c.key)
			require.Error(t, err, "No record error")
			require.Equal(t, "", temp)
			continue
		}
		err := m.Set(c.key, c.shortURL, c.longURL)
		if name == "setError" {
			require.Error(t, err, "record already exists")
			continue
		}
		require.NoError(t, err, name)
		temp, err = m.Get(c.key)
		require.NoError(t, err, name)
		require.Equal(t, c.longURL, temp, name)
	}
}
