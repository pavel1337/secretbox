package main

import (
	"context"
	"flag"
	"io/ioutil"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/pavel1337/secretbox/pkg/crypt"
	rs "github.com/pavel1337/secretbox/pkg/storage/redis"
	"github.com/pavel1337/secretbox/pkg/web"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Redis      string `json:"redis"`
	Encryption string `json:"encryption"`
	Url        string `json:"url"`
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	path := flag.String("c", "config.yml", "Path to a config file")
	secret := flag.String("secret", "xii5ooph2qua2woo4oohahNain2iofie", "Secret key for cookie encryption")

	flag.Parse()

	c, err := parseConfig(*path)
	if err != nil {
		log.Fatalln(err)
	}

	db, err := initRedisClient(c.Redis, 1)
	if err != nil {
		log.Fatalln(err)
	}
	store := rs.NewRedisStore(db)

	crypter, err := crypt.NewAESGCM(c.Encryption)
	if err != nil {
		log.Fatalln(err)
	}

	err = web.New(*secret, store, crypter).Start(*addr)
	log.Fatalln(err)

}

func initRedisClient(addr string, db int) (*redis.Client, error) {
	rc := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:       db,
	})
	err := rc.Ping(context.Background()).Err()
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
