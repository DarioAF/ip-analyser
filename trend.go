package main

// Trends will be managed by country as trend-{ISO}
// Each IP will be a member of it, containing the #invocations

func updateTrend(db DBInterface, user User) {
	key := "trend-" + user.ISOCountry
	member := user.IP

	db.IncrTrend(key, member)
}

func ipScore(db DBInterface, country, ip string) int {
	key := "trend-" + country
	return int(db.RetrieveScore(key, ip))
}

func countryBestScore(db DBInterface, country string) int {
	key := "trend-" + country
	return int(db.TopScore(key).Score)
}

func countryAvgRequests(db DBInterface, country string) int {
	key := "trend-" + country
	countryScores := db.RetrieveAllScores(key)
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
