package service

import (
	"github.com/DarioAF/ip-analyser/pkg/db"
	"github.com/DarioAF/ip-analyser/pkg/model"
)

// Scores will be stored by country as scores-{ISO}
// Each IP will be a member of it, containing the #invocations
var generateKey = func(iso string) string { return "scores-" + iso }

// UpdateScore adds one to the current ip score or creates a new one with score: 1
func UpdateScore(database db.Interface, user model.User) {
	database.IncrScore(generateKey(user.ISOCountry), user.IP)
}

// CountryBestScore returns the top score for the specified country
func CountryBestScore(database db.Interface, country string) int {
	return int(database.TopScore(generateKey(country)).Score)
}

// CountryAvgRequests sums all country scores and divide it by the total members (ip) of that country
func CountryAvgRequests(database db.Interface, country string) int {
	countryScores := database.RetrieveAllScores(generateKey(country))
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
