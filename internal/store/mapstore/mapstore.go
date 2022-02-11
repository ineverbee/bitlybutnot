package mapstore

import (
	"errors"
	"sync"

	"github.com/ineverbee/bitlybutnot/internal/store"
)

type LinkMap struct {
	sync.RWMutex
	mp map[uint32]struct{ short, long string }
}

func (m *LinkMap) Set(key uint32, shortURL, longURL string) error {
	defer m.Unlock()
	m.Lock()
	if _, ok := m.mp[key]; ok {
		return errors.New("record already exists")
	}
	m.mp[key] = struct{ short, long string }{shortURL, longURL}
	return nil
}

func (m *LinkMap) Get(key uint32) (string, error) {
	defer m.RUnlock()
	m.RLock()
	if value, ok := m.mp[key]; ok {
		return value.long, nil
	}
	return "", errors.New("no record found")
}

func NewLinkMap(cap int) store.Store {
	m := &LinkMap{
		mp: make(map[uint32]struct{ short, long string }, cap),
	}
	return m
}
