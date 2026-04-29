package redis

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var counts int

var Ctx = context.Background()

func New(addr string) (*redis.Client, error) {
	var client *redis.Client
	var err error

	maxRetries := 5
	backoff := time.Second

	for {
		client = redis.NewClient(&redis.Options{
			Addr: addr,
		})
		// verify connection
		_, err = client.Ping(Ctx).Result()
		if err != nil {
			log.Println("Redis:conn failed")
			counts++
		} else {
			log.Println("Redis:conn successful")
			return client, nil
		}

		if counts > maxRetries {
			log.Println(err)
			return nil, errors.New("failed to connect to redis after retries: " + err.Error())
		}

		time.Sleep(backoff)
		backoff *= 2 // exponential backoff
		continue
	}
}
