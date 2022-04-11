package main

import (
	"context"
	"database/sql"
	"embed"
	"flag"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/pavel1337/secretbox/pkg/crypt"
	"github.com/pavel1337/secretbox/pkg/storage"
	"github.com/pavel1337/secretbox/pkg/storage/inmem"
	"github.com/pavel1337/secretbox/pkg/storage/mysql"
	rs "github.com/pavel1337/secretbox/pkg/storage/redis"
	"github.com/pavel1337/secretbox/pkg/web"
	"github.com/rbcervilla/redisstore/v8"
)

var (
	addr                 string = os.Getenv("LISTEN_ADDRESS")
	cookieSecret         string = os.Getenv("COOKIE_SECRET")
	sessionStoreType     string = os.Getenv("SESSION_STORE_TYPE")
	secretsStoreType     string = os.Getenv("SECRETS_STORE_TYPE")
	secretsEncryptionKey string = os.Getenv("SECRETS_ENCRYPTION_KEY")
	redisAddr            string = os.Getenv("REDIS_ADDR")
	mysqlDSN             string = os.Getenv("MYSQL_DSN")
	maxCookieAge         int    = 12 * 60 * 60 // 12 hours
)

//go:embed pkg/storage/mysql/migrations/*.sql
var migrationsFS embed.FS

func init() {
	flag.StringVar(&addr, "addr", addr, "HTTP network address")
	flag.StringVar(&cookieSecret, "cookie-key", cookieSecret, "key for cookie encryption")
	flag.StringVar(&secretsEncryptionKey, "secrets-key", secretsEncryptionKey, "key for secrets encryption")
	flag.StringVar(&sessionStoreType, "session-store-type", sessionStoreType, "type of cookie store (REDIS/INMEM(default)")
	flag.StringVar(&secretsStoreType, "secrets-store-type", secretsStoreType, "type of secrets store (REDIS/MYSQL/INMEM(default)")
	flag.StringVar(&redisAddr, "redis-addr", redisAddr, "redis address for redis store (defaults to 127.0.0.1:6379)")
	flag.StringVar(&mysqlDSN, "mysql-dsn", mysqlDSN, "mysql connection string")
}

func main() {
	flag.Parse()
	crypter, err := crypt.NewAESGCM(secretsEncryptionKey)
	if err != nil {
		log.Fatalf("cannot create encryption service due to: %s", err)
	}

	session := cookieStore(sessionStoreType)
	store := secretsStore(secretsStoreType)

	err = web.New(session, store, crypter).Start(addr)
	log.Fatalf("server failed due to: %s", err)

}

func cookieStore(typ string) sessions.Store {
	switch typ {
	case "REDIS":
		db, err := initRedisClient(redisAddr)
		if err != nil {
			log.Fatalf("cannot connect to redis due to: %s", err)
		}
		store, err := redisstore.NewRedisStore(context.Background(), db)
		if err != nil {
			log.Fatalf("cannot create new redis session store due to: %s", err)
		}

		store.Options(sessions.Options{
			MaxAge: maxCookieAge,
		})
		return store

	default:
		var cs []byte
		if cookieSecret != "" {
			cs = []byte(cookieSecret)
		} else {
			cs = securecookie.GenerateRandomKey(32)
		}
		session := sessions.NewCookieStore(cs)
		session.MaxAge(maxCookieAge)
		return session
	}
}

func secretsStore(typ string) storage.Store {
	switch typ {
	case "MYSQL":
		db, err := sql.Open("mysql", mysqlDSN)
		if err != nil {
			log.Fatalf("cannot connect to mysql due to: %s", err)
		}
		db.SetMaxOpenConns(100)
		storage, err := mysql.NewMySQLStorage(db, migrationsFS, "pkg/storage/mysql/migrations")
		if err != nil {
			log.Fatalf("could not create mysql store: %s", err)
		}
		return storage
	case "REDIS":
		db, err := initRedisClient(redisAddr)
		if err != nil {
			log.Fatalf("cannot connect to redis due to: %s", err)
		}
		return rs.NewRedisStore(db)
	default:
		return inmem.NewInmemStore()
	}
}

func initRedisClient(addr string) (*redis.Client, error) {
	rc := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:       0,
	})
	err := rc.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}
	return rc, nil
}
