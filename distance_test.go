package main

import (
	"testing"

	redis "github.com/go-redis/redis/v8"
)

func TestCalcDistance(t *testing.T) {
	distInKm := int(calcDistance(-34.0, -64.0, 64.0, 26.0, "K"))
	if distInKm != 13361 {
		t.Errorf("from Argentina to Finland [latlng (64.0, 26.0)] there are 13361 km and we calculed: %d km", distInKm)
	}
}

func TestResolveDistance(t *testing.T) {
	var db DBInterface = &mockDB{}
	var countryAR IP2countryResponse = IP2countryResponse{"AR", "ARG", "Argentina"}
	var countryFI IP2countryResponse = IP2countryResponse{"FI", "FIN", "Finland"}
	var countryEH IP2countryResponse = IP2countryResponse{"EH", "ESH", "Western Sahara"}

	dist := resolveDistance(db, countryAR)
	if dist != 0 {
		t.Errorf("from Argentina to Argentina there are 0 km and we calculed: %d km", dist)
	}

	dist = resolveDistance(db, countryFI)
	if dist != 13361 {
		t.Errorf("from Argentina to Finland there are 13361 km and we calculed: %d km", dist)
	}

	dist = resolveDistance(db, countryEH)
	if dist != 8444 {
		t.Errorf("from Argentina to Argentina there are 0 km and we calculed: %d km", dist)
	}
}

type mockDB struct{}

func (c *mockDB) Exists(hash, key string) bool {
	if key == "FIN" {
		return true
	}
	return false
}
func (c *mockDB) Retrieve(hash, key string) string {
	return "13361"
}

// Unused methods in this scenario
func (c *mockDB) TopScore(key string) redis.Z              { return redis.Z{} }
func (c *mockDB) RetrieveScore(key, member string) float64 { return 0 }
func (c *mockDB) RetrieveAllScores(key string) []redis.Z   { return []redis.Z{} }
func (c *mockDB) Store(hash, key, value string)            {}
func (c *mockDB) IncrTrend(key, member string)             {}
func (c *mockDB) Ping() string                             { return "PONG" }
