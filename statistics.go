package main

import (
	"log"
	"strconv"
	"strings"
)

// Statistic response for service
type Statistic struct {
	Country  string `json:"country"`
	Distance int    `json:"distance"`
}

func makeStatisticResponse(str string) Statistic {
	info := strings.Split(str, "-")
	dst, err := strconv.Atoi(info[1])
	if err != nil {
		log.Panicf("Cannot convert %s to int", info[1])
	}
	return Statistic{info[0], dst}
}

func updateStat(stat, hash string, user User, strategy func(int, int) bool) {
	if exists(hash, stat) {
		res := retrieve(hash, stat)
		val := strings.Split(res, "-")

		currentStatDistance, err := strconv.Atoi(val[1])
		if err != nil {
			log.Panicf("Cannot convert %s to int", val)
		}

		if !(stat == "nearest" && currentStatDistance == 0) { // If nearest is 0 we're done
			if strategy(currentStatDistance, user.Distance) {
				log.Printf("Storing new %s distance: %s with %d km from AR", stat, user.ISOCountry, user.Distance)
				store(hash, stat, user.ISOCountry+"-"+strconv.Itoa(user.Distance))

			} else if currentStatDistance == user.Distance { // When equals take the one with highest score
				userIPCountryScore := countryBestScore(user.ISOCountry)
				currentStatCountryScore := countryBestScore(val[0])

				if userIPCountryScore > currentStatCountryScore {
					log.Printf("Storing new %s country: %s because greater score. %s has a score of: %d while %s had %d", stat, user.ISOCountry, user.ISOCountry, userIPCountryScore, val[0], currentStatCountryScore)
					store(hash, stat, user.ISOCountry+"-"+strconv.Itoa(user.Distance))
				}
			}
		}

	} else {
		log.Printf("Storing first %s distance: %s with %d km from AR", stat, user.ISOCountry, user.Distance)
		store(hash, stat, user.ISOCountry+"-"+strconv.Itoa(user.Distance))
	}
}

func farthestStrategy(current, actual int) bool {
	return current < actual
}
func nearestStrategy(current, actual int) bool {
	return current > actual
}

func updateStatistics(user User) {
	fromIata := "AR"
	hash := "statistics-" + fromIata

	updateStat("farthest", hash, user, farthestStrategy)
	updateStat("nearest", hash, user, nearestStrategy)
}

func retrieveFarthest() Statistic {
	fromIata := "AR"
	hash := "statistics-" + fromIata
	key := "farthest"

	return makeStatisticResponse(retrieve(hash, key))
}

func retrieveNearest() Statistic {
	fromIata := "AR"
	hash := "statistics-" + fromIata
	key := "nearest"

	return makeStatisticResponse(retrieve(hash, key))
}
