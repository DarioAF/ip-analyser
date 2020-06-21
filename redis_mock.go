package main

import redis "github.com/go-redis/redis/v8"

type mockDB struct {
	ping              string
	exists            func(hash, key string) bool
	retrieve          func(hash, key string) string
	topScore          func(key string) redis.Z
	retrieveScore     func(key, member string) float64
	retrieveAllScores func(key string) []redis.Z
}

func (c *mockDB) Exists(hash, key string) bool {
	return c.exists(hash, key)
}

func (c *mockDB) Retrieve(hash, key string) string {
	return c.retrieve(hash, key)
}

func (c *mockDB) TopScore(key string) redis.Z {
	return c.topScore(key)
}

func (c *mockDB) RetrieveScore(key, member string) float64 {
	return c.retrieveScore(key, member)
}

func (c *mockDB) RetrieveAllScores(key string) []redis.Z {
	return c.retrieveAllScores(key)
}

func (c *mockDB) Ping() string {
	return c.ping
}

func (c *mockDB) Store(hash, key, value string) {}
func (c *mockDB) IncrTrend(key, member string)  {}
