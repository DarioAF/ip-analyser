package main

// Trends will be managed by country as trend-ISO
// Each IP will be a member of it, containing the #invocations

func updateTrend(user User) {
	key := "trend-" + user.ISOCountry
	member := user.IP

	incrTrend(key, member)
}

func ipScore(country, ip string) int {
	key := "trend-" + country
	return int(retrieveScore(key, ip))
}

func countryBestScore(country string) int {
	key := "trend-" + country
	return int(topScore(key).Score)
}

func countryAvgRequests(country string) int {
	key := "trend-" + country
	countryScores := retrieveAllScores(key)
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
