package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/pavel1337/secretbox/pkg/storage"
)

type redisStore struct {
	rdb *redis.Client
}

var ctx = context.Background()

func NewRedisStore(client *redis.Client) *redisStore {
	return &redisStore{client}
}

func (ss *redisStore) Insert(s storage.Secret, ttlMinutes int) (string, error) {
	id := uuid.New().String()
	bb, err := json.Marshal(s)
	if err != nil {
		return id, err
	}

	ttlDur := (time.Duration(ttlMinutes) * time.Minute)
	_, err = ss.rdb.Set(ctx, id, bb, ttlDur).Result()
	if err != nil {
		return id, err
	}

	return id, err
}

func (ss *redisStore) Exists(id string) bool {
	_, err := ss.rdb.Get(ctx, id).Result()
	if err != nil {
		return false
	}

	return true
}

func (ss *redisStore) Get(id string) (*storage.Secret, error) {
	bb, err := ss.rdb.Get(ctx, id).Bytes()
	if err == redis.Nil {
		return nil, storage.ErrNoRecord
	}
	if err != nil {
		return nil, err
	}

	var s storage.Secret

	err = json.Unmarshal(bb, &s)
	if err != nil {
		return nil, err
	}

	s.ID = id

	return &s, nil
}

func (ss *redisStore) Delete(key string) error {
	i, err := ss.rdb.Del(ctx, key).Result()
	if err != nil {
		return err
	}
	if i == 0 {
		return storage.ErrNoRecord
	}
	return nil
}
