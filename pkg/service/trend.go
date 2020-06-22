package service

import (
	"github.com/DarioAF/ip-analyser/pkg/db"
	"github.com/DarioAF/ip-analyser/pkg/model"
)

// Trends will be managed by country as trend-{ISO}
// Each IP will be a member of it, containing the #invocations

func UpdateTrend(database db.DBInterface, user model.User) {
	key := "trend-" + user.ISOCountry
	member := user.IP

	database.IncrTrend(key, member)
}

func IpScore(database db.DBInterface, country, ip string) int {
	key := "trend-" + country
	return int(database.RetrieveScore(key, ip))
}

func CountryBestScore(database db.DBInterface, country string) int {
	key := "trend-" + country
	return int(database.TopScore(key).Score)
}

func CountryAvgRequests(database db.DBInterface, country string) int {
	key := "trend-" + country
	countryScores := database.RetrieveAllScores(key)
	members := len(countryScores)
	if members == 0 {
		return 0
	}
	countryScoreSum := 0
	for _, n := range countryScores {
		countryScoreSum += int(n.Score)
	}
	return countryScoreSum / members
}
