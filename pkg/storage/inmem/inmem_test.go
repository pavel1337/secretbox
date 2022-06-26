package inmem

import (
	"testing"

	"github.com/pavel1337/secretbox/pkg/storage"
	"github.com/stretchr/testify/assert"
)

var testingContent = "testing-content"

var secret = storage.Secret{
	Content:    []byte(testingContent),
	Passphrase: "passphrase",
}

func TestInmemStore(t *testing.T) {
	s := NewInmemStore()

	id, err := s.Insert(secret, 60)
	assert.NoError(t, err)
	assert.True(t, s.Exists(id))

	newS, err := s.GetAndDelete(id)
	assert.NoError(t, err)

	assert.Equal(t, newS.Passphrase, secret.Passphrase)
	assert.Equal(t, newS.ID, id)
	assert.Equal(t, string(newS.Content), testingContent)

	s.Delete(id)
	assert.NoError(t, err)

	assert.False(t, s.Exists(id))

	newS, err = s.GetAndDelete(id)
	assert.Error(t, err)
	assert.ErrorIs(t, err, storage.ErrNoRecord)

	id, _ = s.Insert(secret, -60)
	assert.False(t, s.Exists(id))
	newS, err = s.GetAndDelete(id)
	assert.Error(t, err)
	assert.ErrorIs(t, err, storage.ErrNoRecord)

}
