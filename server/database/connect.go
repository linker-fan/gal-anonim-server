package database

import (
	"database/sql"
	"fmt"
	"linker-fan/gal-anonim-server/server/config"
	"log"

	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
)

var db *sql.DB

func ConnectToPostgres(c *config.Config) error {
	psqlInfo := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", c.Postgres.User, c.Postgres.Password, c.Postgres.Host, c.Postgres.Port, c.Postgres.Name)
	log.Println(psqlInfo)

	database, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("sql.Open failed: %v\n", err)
	}

	err = database.Ping()
	if err != nil {
		log.Fatalf("database.Ping failed: %v\n", err)
	}

	db = database
	log.Println("[+] Connected to the database")
	return nil
}

func ConnectToRedis(c *config.Config) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", c.Redis.Host, c.Redis.Password),
		Password: c.Redis.Password, // no password set
		DB:       c.Redis.DB,       // use default DB
	})

	return rdb
}
