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
		log.Printf("ERROR: cannot convert %s to int", info[1])
	}
	return Statistic{info[0], dst}
}

func updateStat(db DBInterface, stat, hash string, user User, strategy func(int, int) bool) {
	if db.Exists(hash, stat) {
		res := db.Retrieve(hash, stat)
		val := strings.Split(res, "-")

		currentStatDistance, err := strconv.Atoi(val[1])
		if err != nil {
			log.Printf("ERROR: cannot convert %s to int", val)
		}

		if !(stat == "nearest" && currentStatDistance == 0) { // If nearest is 0 we're done
			if strategy(currentStatDistance, user.Distance) {
				log.Printf("Storing new %s distance: %s with %d km from AR", stat, user.ISOCountry, user.Distance)
				db.Store(hash, stat, user.ISOCountry+"-"+strconv.Itoa(user.Distance))

			} else if currentStatDistance == user.Distance { // When equals take the one with (country) highest score
				userIPCountryScore := countryBestScore(db, user.ISOCountry)
				currentStatCountryScore := countryBestScore(db, val[0])

				if userIPCountryScore > currentStatCountryScore {
					log.Printf("Storing new %s country: %s because greater score. %s has a score of: %d while %s had %d", stat, user.ISOCountry, user.ISOCountry, userIPCountryScore, val[0], currentStatCountryScore)
					db.Store(hash, stat, user.ISOCountry+"-"+strconv.Itoa(user.Distance))
				}
			}
		}

	} else {
		log.Printf("Storing first %s distance: %s with %d km from AR", stat, user.ISOCountry, user.Distance)
		db.Store(hash, stat, user.ISOCountry+"-"+strconv.Itoa(user.Distance))
	}
}

func farthestStrategy(current, actual int) bool {
	return current < actual
}
func nearestStrategy(current, actual int) bool {
	return current > actual
}

func updateStatistics(db DBInterface, user User) {
	fromIata := "AR"
	hash := "statistics-" + fromIata

	updateStat(db, "farthest", hash, user, farthestStrategy)
	updateStat(db, "nearest", hash, user, nearestStrategy)
}

func retrieveDistance(db DBInterface, key string) Statistic {
	fromIata := "AR"
	hash := "statistics-" + fromIata

	return makeStatisticResponse(db.Retrieve(hash, key))
}
