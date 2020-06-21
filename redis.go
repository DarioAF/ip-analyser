package main

import (
	"log"
	"os"

	redis "github.com/go-redis/redis/v8"
)

// DBInterface will provide us the expected methods or db must have
type DBInterface interface {
	Ping() string
	Exists(hash, key string) bool
	Retrieve(hash, key string) string
	Store(hash, key, value string)
	IncrTrend(key, member string)
	TopScore(key string) redis.Z
	RetrieveScore(key, member string) float64
	RetrieveAllScores(key string) []redis.Z
}

// RedisClient is the redis abstraction of the defined db interface
type RedisClient struct {
	client *redis.Client
}

// RedisConnector is the redis implementation of the defined RedisClient
var RedisConnector = RedisClient{
	redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})}

// Ping call will return PONG if the connection was successful
func (c *RedisClient) Ping() string {
	res, err := c.client.Ping(c.client.Context()).Result()
	if err != nil {
		log.Print("ERROR: there was an error connecting with redis")
	}
	return res
}

// Exists returns true if the key inside that hash exists, or false if not
func (c *RedisClient) Exists(hash, key string) bool {
	res, err := c.client.HExists(c.client.Context(), hash, key).Result()
	if err != nil {
		log.Printf("ERROR: there was an error trying to check %s existence in %s", key, hash)
	}
	return res
}

// Retrieve will get the value for that key in the given hash
func (c *RedisClient) Retrieve(hash, key string) string {
	res, err := c.client.HGet(c.client.Context(), hash, key).Result()
	if err != nil {
		log.Printf("ERROR: there was an error trying to retrieve %s from %s", key, hash)
	}
	return res
}

// Store sets a key -> value inside a hash
func (c *RedisClient) Store(hash, key, value string) {
	_, err := c.client.HSet(c.client.Context(), hash, key, value).Result()
	if err != nil {
		log.Printf("ERROR: there was an error storing %s into %s for %s", value, key, hash)
	}
}

// IncrTrend will add one to the total score for the given member of that key
func (c *RedisClient) IncrTrend(key, member string) {
	_, err := c.client.ZIncrBy(c.client.Context(), key, 1, member).Result()
	if err != nil {
		log.Printf("ERROR: there was an error updating %s for %s trend", key, member)
	}
}

// TopScore will return the (member -> score) pair with highest score for the given key
func (c *RedisClient) TopScore(key string) redis.Z {
	res, err := c.client.ZRangeWithScores(c.client.Context(), key, -1, -1).Result()
	if err != nil {
		log.Printf("ERROR: there was an error retrieving %s trend", key)
	}
	return res[0]
}

// RetrieveScore will, return the score for the given member for that key
func (c *RedisClient) RetrieveScore(key, member string) float64 {
	res, err := c.client.ZScore(c.client.Context(), key, member).Result()
	if err != nil {
		log.Printf("ERROR: there was an error retrieving score for ip %s in %s", member, key)
	}
	return res
}

// RetrieveAllScores will return all the (member -> score) pairs for the given key
func (c *RedisClient) RetrieveAllScores(key string) []redis.Z {
	res, err := c.client.ZRangeWithScores(c.client.Context(), key, 0, -1).Result()
	if err != nil {
		log.Printf("ERROR: there was an error retrieving %s trend", key)
	}
	return res
}
