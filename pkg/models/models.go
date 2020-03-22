package models

import (
	"errors"
	"time"

	"github.com/go-redis/redis"
	"github.com/google/uuid"
)

type Secret struct {
	Key        string
	EncContent []byte
	Content    string
}

type SecretModel struct {
	DB *redis.Client
}

var ErrNoRecord = errors.New("models: no matching record found")

func (m *SecretModel) Insert(encContent []byte, expires int) (string, error) {
	key := uuid.New().String()
	_, err := m.DB.Set(key, encContent, (time.Duration(expires) * time.Minute)).Result()
	if err != nil {
		return key, err
	}
	return key, err
}

func (m *SecretModel) Get(key string) (*Secret, error) {
	i, err := m.DB.Exists(key).Result()
	if err != nil {
		return nil, err
	}
	if i == 0 {
		return nil, ErrNoRecord
	}
	encContent, err := m.DB.Get(key).Bytes()
	if err != nil {
		return nil, err
	}
	return &Secret{Key: key, EncContent: encContent}, nil
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
