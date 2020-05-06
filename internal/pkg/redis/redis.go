package redis

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/go-redis/redis"
)

var (
	host   = "redis"
	port   = "6379"
	db     = "0"
	once   sync.Once
	client *redis.Client
)

// Get returns a specified number of random tokens from a set of a part of speech.
func Get(set string, number int64) ([]string, error) {
	c := getClient()

	t, err := c.SRandMemberN(set, number).Result()
	if err != nil {
		return nil, err
	}

	return t, nil
}

// Save adds one or more tokens to the appropriate set for a part of speech.
func Save(set string, tokens []string) error {
	c := getClient()

	err := c.SAdd(set, tokens).Err()
	if err != nil {
		return err
	}

	return nil
}

func getClient() *redis.Client {
	once.Do(func() {
		db, err := strconv.Atoi(db)
		if err != nil {
			panic(err)
		}

		client = redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("%s:%s", host, port),
			DB:   db,
		})
	})
	return client
}

func init() {
	hostEnv := os.Getenv("REDIS_HOST")
	if hostEnv != "" {
		host = hostEnv
	}

	portEnv := os.Getenv("REDIS_PORT")
	if portEnv != "" {
		host = portEnv
	}

	dbEnv := os.Getenv("REDIS_DB")
	if dbEnv != "" {
		host = dbEnv
	}
}
