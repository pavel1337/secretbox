package inmem

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/pavel1337/secretbox/pkg/storage"
)

type inmemStore struct {
	m    map[string]storage.Secret
	mTtl map[string]time.Time
	lock sync.RWMutex
}

func NewInmemStore() *inmemStore {
	m := make(map[string]storage.Secret)
	mTtl := make(map[string]time.Time)
	is := &inmemStore{
		m:    m,
		mTtl: mTtl,
		lock: sync.RWMutex{},
	}

	return is
}

func (is *inmemStore) Insert(s storage.Secret, ttlMinutes int) (string, error) {
	id := uuid.New().String()

	is.lock.Lock()
	defer is.lock.Unlock()

	is.m[id] = s
	is.mTtl[id] = time.Now().Add(time.Duration(ttlMinutes) * time.Minute)

	return id, nil
}

func (is *inmemStore) Exists(id string) bool {
	is.lock.RLock()
	defer is.lock.RUnlock()
	_, ok := is.m[id]
	if !ok {
		return false
	}
	t, ok := is.mTtl[id]
	if !ok {
		return false
	}
	if t.Before(time.Now()) {
		return false
	}
	return true

}

func (is *inmemStore) Get(id string) (*storage.Secret, error) {
	is.lock.RLock()
	defer is.lock.RUnlock()
	t, ok := is.mTtl[id]
	if !ok || t.Before(time.Now()) {
		return nil, storage.ErrNoRecord
	}
	secret, ok := is.m[id]
	if !ok {
		return nil, storage.ErrNoRecord
	}
	secret.ID = id
	return &secret, nil
}

func (is *inmemStore) Delete(id string) error {
	is.lock.Lock()
	defer is.lock.Unlock()

	_, ok := is.m[id]
	if !ok {
		return storage.ErrNoRecord
	}

	delete(is.m, id)
	delete(is.mTtl, id)

	return nil
}
