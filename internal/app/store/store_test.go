package store

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func newTestStore(t *testing.T) *Store {
	t.Helper()
	return New()
}

func TestStore_Open(t *testing.T) {
	s := newTestStore(t)
	err := s.Open(dbURL)
	assert.NoError(t, err)
}

func TestStore_Close(t *testing.T) {
	s := newTestStore(t)
	defer func() {
		if r := recover(); r != nil {
			assert.NotNil(t, r)
		}
	}()
	s.Close()
}
