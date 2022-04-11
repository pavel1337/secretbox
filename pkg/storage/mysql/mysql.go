package mysql

import (
	"context"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io/fs"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/pavel1337/secretbox/pkg/storage"
	"github.com/pressly/goose/v3"
)

func NewMySQLStorage(db *sql.DB, migrationsFS fs.FS, migrationsDir string) (*Queries, error) {
	q := New(db)
	goose.SetBaseFS(migrationsFS)

	if err := goose.SetDialect("mysql"); err != nil {
		return nil, fmt.Errorf("could not set dialect: %w", err)
	}

	if err := goose.Up(db, migrationsDir); err != nil {
		return nil, fmt.Errorf("could not run migrations: %w", err)
	}

	go q.deleteExpired()

	return q, nil
}

func (q *Queries) deleteExpired() {
	ticker := time.NewTicker(10 * time.Second)
	for range ticker.C {
		ids, err := q.listExpiredSecretsIds(context.Background(), time.Now())
		if err != nil {
			log.Printf("could not list expired secrets: %s", err)
			continue
		}
		for _, id := range ids {
			err = q.Delete(id)
			if err != nil {
				log.Printf("could not delete expired secret %s due to %s", id, err)
			}
		}

	}
}

func (q *Queries) Insert(s storage.Secret, ttlMinutes int) (string, error) {
	id := uuid.New().String()

	p := createSecretParams{
		ID:         id,
		Content:    hex.EncodeToString(s.Content),
		Passphrase: sql.NullString{String: s.Passphrase, Valid: true},
		ExpiresAt:  time.Now().Add(time.Duration(ttlMinutes) * time.Minute),
	}

	_, err := q.createSecret(context.Background(), p)
	if err != nil {
		return id, err
	}

	return id, nil
}

func (q *Queries) Exists(id string) bool {
	_, err := q.Get(id)
	if err != nil {
		return false
	}
	return true
}

func (q *Queries) Get(id string) (*storage.Secret, error) {
	s, err := q.getSecret(context.Background(), id)
	if err == sql.ErrNoRows {
		return nil, storage.ErrNoRecord
	}
	if err != nil {
		return nil, err
	}
	if s.ExpiresAt.Before(time.Now()) {
		return nil, storage.ErrNoRecord
	}

	sec, err := dehexifySecret(s)
	if err != nil {
		return nil, err
	}

	return sec, nil
}

func (q *Queries) Delete(id string) error {
	err := q.deleteSecret(context.Background(), id)
	if err == sql.ErrNoRows {
		return storage.ErrNoRecord
	}
	if err != nil {
		return err
	}
	return nil
}

func dehexifySecret(s Secret) (*storage.Secret, error) {
	content, err := hex.DecodeString(s.Content)
	if err != nil {
		return nil, fmt.Errorf("could not decode string:%w", err)
	}

	var passphrase string
	if s.Passphrase.Valid {
		passphrase = s.Passphrase.String
	}

	return &storage.Secret{
		ID:         s.ID,
		Content:    content,
		Passphrase: passphrase,
	}, nil
}
