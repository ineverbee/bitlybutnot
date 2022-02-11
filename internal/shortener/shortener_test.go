package shortener

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHashURL(t *testing.T) {
	cases := map[string]struct {
		url      string
		expected uint32
	}{
		"example":      {"example.com", 1125968678},
		"youtube":      {"https://www.youtube.com/", 4083526518},
		"google":       {"https://www.google.ru/", 2250201066},
		"some link":    {"https://www.youtube.com/watch?v=dQw4w9WgXcQ", 4255220111},
		"e":            {"e", 3758891744},
		"empty string": {"", 2166136261},
	}

	for name, c := range cases {
		require.Equal(t, c.expected, HashURL(c.url), name)
	}
}

func TestUniqueString(t *testing.T) {
	cases := map[string]struct {
		h        uint32
		expected string
	}{
		"example":      {1125968678, "SWbEibaaaa"},
		"youtube":      {4083526518, "pqcoheaaaa"},
		"google":       {2250201066, "pfh2qcaaaa"},
		"some link":    {4255220111, "x6Qhseaaaa"},
		"e":            {3758891744, "uQUMXdaaaa"},
		"empty string": {2166136261, "qU6Flcaaaa"},
	}

	for name, c := range cases {
		require.Equal(t, c.expected, UniqueString(c.h), name)
	}
}

func TestDecodeString(t *testing.T) {
	cases := map[string]struct {
		link     string
		expected uint32
	}{
		"example":      {"SWbEibaaaa", 1125968678},
		"youtube":      {"pqcoheaaaa", 4083526518},
		"google":       {"pfh2qcaaaa", 2250201066},
		"some link":    {"x6Qhseaaaa", 4255220111},
		"e":            {"uQUMXdaaaa", 3758891744},
		"empty string": {"qU6Flcaaaa", 2166136261},
		"error":        {"~jndUaaaaa", 0},
	}

	for name, c := range cases {
		l, err := DecodeString(c.link)
		if name == "error" {
			require.Error(t, err, err)
			continue
		}
		require.Nil(t, err)
		require.Equal(t, c.expected, l, name)
	}
}

func BenchmarkHashURL(b *testing.B) {
	for n := 0; n < b.N; n++ {
		HashURL("example.com")
	}
}

func BenchmarkUniqueString(b *testing.B) {
	for n := 0; n < b.N; n++ {
		UniqueString(2166136261)
	}
}

func BenchmarkDecodeString(b *testing.B) {
	for n := 0; n < b.N; n++ {
		if n%2 == 0 {
			DecodeString("D2Bpvaaaaa")
		} else {
			DecodeString("~jndUaaaaa")
		}
	}
}
