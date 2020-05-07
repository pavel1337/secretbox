package models

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"github.com/pavel1337/secretbox/pkg/crypt"
)

type Secret struct {
	Key        string `json:"key"`
	Content    string `json:"content"`
	Passphrase string `json:"passphrase"`
}

type SecretModel struct {
	DB *redis.Client
}

var ErrNoRecord = errors.New("models: no matching record found")
var ErrIncorrectPassphrase = errors.New("models: incorrect assphrase")

func (m *SecretModel) Insert(s Secret, e *[32]byte, expires int) (string, error) {
	key := uuid.New().String()
	j, err := json.Marshal(s)
	if err != nil {
		return key, err
	}
	value, err := crypt.EncryptSecret(j, e)
	if err != nil {
		return key, err
	}
	_, err = m.DB.Set(key, value, (time.Duration(expires) * time.Minute)).Result()
	if err != nil {
		return key, err
	}
	return key, err
}

func (m *SecretModel) Exists(key string) error {
	i, err := m.DB.Exists(key).Result()
	if err != nil {
		return err
	}
	if i == 0 {
		return ErrNoRecord
	}
	return nil
}

func (m *SecretModel) Get(key string, e *[32]byte) (*Secret, error) {
	i, err := m.DB.Exists(key).Result()
	if err != nil {
		return nil, err
	}
	if i == 0 {
		return nil, ErrNoRecord
	}
	value, err := m.DB.Get(key).Bytes()
	if err != nil {
		return nil, err
	}
	j, err := crypt.DecryptSecret(value, e)
	if err != nil {
		return nil, err
	}
	var s Secret
	err = json.Unmarshal(j, &s)
	if err != nil {
		return nil, err
	}
	s.Key = key
	return &s, nil
}

func (m *SecretModel) Delete(key string) error {
	i, err := m.DB.Del(key).Result()
	if err != nil {
		return err
	}
	if i == 0 {
		return ErrNoRecord
	}
	return nil
}
