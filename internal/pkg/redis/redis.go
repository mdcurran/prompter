package redis

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var (
	host   = "localhost"
	port   = "6379"
	db     = "0"
	once   sync.Once
	client *redis.Client
)

// Ping performs a health check of the Redis connection by pinging the server.
func Ping() error {
	c := getClient()

	err := c.Ping().Err()
	if err != nil {
		return err
	}

	return nil
}

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

// Record saves a sentence to Redis.
func Record(sentence string) error {
	c := getClient()

	err := c.SAdd("sentences", sentence).Err()
	if err != nil {
		return err
	}

	return nil
}

// VerifyUser returns true if the user supplies a password matching the hash stored in Redis.
func VerifyUser(email, password string) bool {
	c := getClient()

	hashed, err := c.HGet("users", email).Result()
	if err != nil {
		return false
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	if err != nil {
		return false
	}

	return true
}

// AddUser adds an email and username entry to Redis.
func AddUser(email, password string) error {
	c := getClient()

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil
	}

	err = c.HSet("users", email, hashed).Err()
	if err != nil {
		return err
	}

	return nil
}

func getClient() *redis.Client {
	once.Do(func() {
		db, err := strconv.Atoi(db)
		if err != nil {
			zap.S().Error(err)
			os.Exit(1)
		}

		zap.S().Info("creating Redis client")
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
