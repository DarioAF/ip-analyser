package db

import redis "github.com/go-redis/redis/v8"

type MockDB struct {
	PingMock              string
	ExistsMock            func(hash, key string) bool
	RetrieveMock          func(hash, key string) string
	TopScoreMock          func(key string) redis.Z
	RetrieveAllScoresMock func(key string) []redis.Z
}

func (c *MockDB) Exists(hash, key string) bool {
	return c.ExistsMock(hash, key)
}

func (c *MockDB) Retrieve(hash, key string) string {
	return c.RetrieveMock(hash, key)
}

func (c *MockDB) TopScore(key string) redis.Z {
	return c.TopScoreMock(key)
}

func (c *MockDB) RetrieveAllScores(key string) []redis.Z {
	return c.RetrieveAllScoresMock(key)
}

func (c *MockDB) Ping() string {
	return c.PingMock
}

func (c *MockDB) Store(hash, key, value string) {}
func (c *MockDB) IncrScore(key, member string)  {}
