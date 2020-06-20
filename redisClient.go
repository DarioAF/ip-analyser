package main

import (
	"log"
	"os"

	redis "github.com/go-redis/redis/v8"
)

var redisClient = newRedisClient()

func newRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func ping() string {
	res, err := redisClient.Ping(redisClient.Context()).Result()
	if err != nil {
		log.Panic("There was an error connecting with redis !")
	}
	return res
}

func exists(hash, key string) bool {
	res, err := redisClient.HExists(redisClient.Context(), hash, key).Result()
	if err != nil {
		log.Panicf("There was an error trying to check %s existence in %s", key, hash)
	}
	return res
}

func retrieve(hash, key string) string {
	res, err := redisClient.HGet(redisClient.Context(), hash, key).Result()
	if err != nil {
		log.Panicf("There was an error trying to retrieve %s from %s", key, hash)
	}
	return res
}

func store(hash, key, value string) {
	_, err := redisClient.HSet(redisClient.Context(), hash, key, value).Result()
	if err != nil {
		log.Panicf("There was an error storing %s into %s for %s", value, key, hash)
	}
}

func incrTrend(key, member string) {
	_, err := redisClient.ZIncrBy(redisClient.Context(), key, 1, member).Result()
	if err != nil {
		log.Panicf("There was an error updating %s for %s trend", key, member)
	}
}

func topScore(key string) redis.Z {
	res, err := redisClient.ZRangeWithScores(redisClient.Context(), key, -1, -1).Result()
	if err != nil {
		log.Panicf("There was an error retrieving %s trend", key)
	}
	return res[0]
}

func retrieveScore(key, member string) float64 {
	res, err := redisClient.ZScore(redisClient.Context(), key, member).Result()
	if err != nil {
		log.Panicf("There was an error retrieving score for ip %s in %s", member, key)
	}
	return res
}

func retrieveAllScores(key string) []redis.Z {
	res, err := redisClient.ZRangeWithScores(redisClient.Context(), key, 0, -1).Result()
	if err != nil {
		log.Panicf("There was an error retrieving %s trend", key)
	}
	return res
}
