package mysql

import (
	"database/sql"
	"os"
	"testing"

	"github.com/pavel1337/secretbox/pkg/storage"
	"github.com/stretchr/testify/assert"

	_ "github.com/go-sql-driver/mysql"
)

var testingContent = "testing-content"

var secret = storage.Secret{
	Content:    []byte(testingContent),
	Passphrase: "passphrase",
}

func TestMysqlStore(t *testing.T) {
	db, err := sql.Open("mysql", os.Getenv("TEST_MYSQL_DSN"))
	assert.NoError(t, err)

	dir := os.DirFS(".")

	s, err := NewMySQLStorage(db, dir, "migrations")
	assert.NoError(t, err)
	id, err := s.Insert(secret, 60)
	assert.NoError(t, err)
	assert.True(t, s.Exists(id))

	newS, err := s.Get(id)
	assert.NoError(t, err)

	assert.Equal(t, newS.Passphrase, secret.Passphrase)
	assert.Equal(t, newS.ID, id)
	assert.Equal(t, string(newS.Content), testingContent)

	s.Delete(id)
	assert.NoError(t, err)

	assert.False(t, s.Exists(id))

	newS, err = s.Get(id)
	assert.Error(t, err)
	assert.ErrorIs(t, err, storage.ErrNoRecord)

	id, _ = s.Insert(secret, -60)
	assert.False(t, s.Exists(id))
	newS, err = s.Get(id)
	assert.Error(t, err)
	assert.ErrorIs(t, err, storage.ErrNoRecord)

}
