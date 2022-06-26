package storage

import "errors"

type Secret struct {
	ID         string `json:"uuid"`
	Content    []byte `json:"content"`
	Passphrase string `json:"passphrase"`
}

type Store interface {
	Insert(s Secret, ttlMinutes int) (string, error)
	Exists(id string) bool
	Get(id string) (*Secret, error)
	GetAndDelete(id string) (*Secret, error)
	Delete(id string) error
}

var ErrNoRecord = errors.New("no matching record found")
