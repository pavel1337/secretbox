package main

import (
	"crypto/sha256"
	"flag"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis"
	"github.com/golangcollege/sessions"
	"github.com/pavel1337/secretbox/pkg/models"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Redis         string `json:"redis"`
	EncryptionKey string `json:"encryption_key"`
	Url           string `json:"url"`
}

type application struct {
	config        Config
	encryptionKey *[32]byte
	errorLog      *log.Logger
	infoLog       *log.Logger
	session       *sessions.Session
	secrets       *models.SecretModel
	templateCache map[string]*template.Template
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	path := flag.String("c", "config.yml", "Path to a config file")
	secret := flag.String("secret", "xii5ooph2qua2woo4oohahNain2iofie", "Secret key for cookie encryption")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	c, err := parseConfig(*path)
	if err != nil {
		errorLog.Fatal(err)
	}

	db, err := initRedisClient(c.Redis, 1)
	if err != nil {
		errorLog.Fatal(err)
	}

	templateCache, err := newTemplateCache("./ui/html")
	if err != nil {
		errorLog.Fatal(err)
	}

	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour

	encryptionKey := sha256.Sum256([]byte(c.EncryptionKey))

	app := &application{
		config:        c,
		encryptionKey: &encryptionKey,
		errorLog:      errorLog,
		infoLog:       infoLog,
		session:       session,
		secrets:       &models.SecretModel{DB: db},
		templateCache: templateCache,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s\n", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func initRedisClient(addr string, db int) (*redis.Client, error) {
	rc := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:       db,
	})
	err := rc.Ping().Err()
	if err != nil {
		return nil, err
	}
	return rc, nil
}

func parseConfig(p string) (Config, error) {
	c := Config{}
	rawConfig, err := ioutil.ReadFile(p)
	if err != nil {
		flag.Usage()
		return c, err
	}
	err = yaml.Unmarshal(rawConfig, &c)
	if err != nil {
		return c, err
	}
	return c, nil
}
