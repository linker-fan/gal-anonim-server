package database

import (
	"context"
	"database/sql"
	"fmt"
	"linker-fan/gal-anonim-server/server/config"
	"log"

	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
)

//global variables
var db *sql.DB
var RedisClient *redis.Client

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

func ConnectToRedis(c *config.Config) error {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", c.Redis.Host, c.Redis.Password),
		Password: c.Redis.Password, // no password set
		DB:       c.Redis.DB,       // use default DB
	})

	//check connection
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Println(err)
		return err
	}

	RedisClient = rdb
	return nil

}
