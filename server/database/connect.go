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
type DatabaseWrapper struct {
	db          *sql.DB
	RedisClient *redis.Client
}

func NewDatabaseWrapper(c *config.Config) (*DatabaseWrapper, error) {
	var dw *DatabaseWrapper
	db, err := connectToPostgres(c)
	if err != nil {
		return nil, err
	}

	rdb, err := connectToRedis(c)
	if err != nil {
		return nil, err
	}

	dw.db = db
	dw.RedisClient = rdb

	return dw, nil
}

func connectToPostgres(c *config.Config) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", c.Postgres.User, c.Postgres.Password, c.Postgres.Host, c.Postgres.Port, c.Postgres.Name)
	log.Println(psqlInfo)

	database, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("sql.Open failed: %v\n", err)
		return nil, err
	}

	err = database.Ping()
	if err != nil {
		log.Fatalf("database.Ping failed: %v\n", err)
		return nil, err
	}

	return database, nil
}

func connectToRedis(c *config.Config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", c.Redis.Host, c.Redis.Port),
		Password: c.Redis.Password, // no password set
		DB:       c.Redis.DB,       // use default DB
	})

	//check connection
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Printf("rdb.Ping failed: %v\n", err)
		return nil, err
	}

	return rdb, nil
}
