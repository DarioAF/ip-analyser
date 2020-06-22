package service

import (
	"log"
	"strconv"
	"strings"

	"github.com/DarioAF/ip-analyser/pkg/db"
	"github.com/DarioAF/ip-analyser/pkg/model"
)

func makeStatisticResponse(str string) model.Statistic {
	info := strings.Split(str, "-")
	dst, err := strconv.Atoi(info[1])
	if err != nil {
		log.Printf("ERROR: cannot convert %s to int", info[1])
	}
	return model.Statistic{info[0], dst}
}

func updateStat(database db.Interface, stat, hash string, user model.User, strategy func(int, int) bool) {
	if database.Exists(hash, stat) {
		res := database.Retrieve(hash, stat)
		val := strings.Split(res, "-")

		currentStatDistance, err := strconv.Atoi(val[1])
		if err != nil {
			log.Printf("ERROR: cannot convert %s to int", val)
		}

		if !(stat == "nearest" && currentStatDistance == 0) { // If nearest is 0 we're done
			if strategy(currentStatDistance, user.Distance) {
				log.Printf("Storing new %s distance: %s with %d km from AR", stat, user.ISOCountry, user.Distance)
				database.Store(hash, stat, user.ISOCountry+"-"+strconv.Itoa(user.Distance))

			} else if currentStatDistance == user.Distance { // When equals take the one with (country) highest score
				userIPCountryScore := CountryBestScore(database, user.ISOCountry)
				currentStatCountryScore := CountryBestScore(database, val[0])

				if userIPCountryScore > currentStatCountryScore {
					log.Printf("Storing new %s country: %s because greater score. %s has a score of: %d while %s had %d", stat, user.ISOCountry, user.ISOCountry, userIPCountryScore, val[0], currentStatCountryScore)
					database.Store(hash, stat, user.ISOCountry+"-"+strconv.Itoa(user.Distance))
				}
			}
		}

	} else {
		log.Printf("Storing first %s distance: %s with %d km from AR", stat, user.ISOCountry, user.Distance)
		database.Store(hash, stat, user.ISOCountry+"-"+strconv.Itoa(user.Distance))
	}
}

func farthestStrategy(current, actual int) bool {
	return current < actual
}
func nearestStrategy(current, actual int) bool {
	return current > actual
}

func UpdateStatistics(database db.Interface, user model.User) {
	fromIata := "AR"
	hash := "statistics-" + fromIata

	updateStat(database, "farthest", hash, user, farthestStrategy)
	updateStat(database, "nearest", hash, user, nearestStrategy)
}

func RetrieveDistance(database db.Interface, key string) model.Statistic {
	fromIata := "AR"
	hash := "statistics-" + fromIata

	return makeStatisticResponse(database.Retrieve(hash, key))
}
