package service

import (
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/DarioAF/ip-analyser/pkg/db"
	"github.com/DarioAF/ip-analyser/pkg/model"
)

// Statistics will be stored with hash: statistics-{(from) country ISO}
// Since its allways from AR, is static
var statisticshash string = "statistics-AR"

// When stored, values will be composed by {(to) country ISO}-{distance (from) country ISO}
var toValue = func(user model.User) string {
	return user.ISOCountry + "-" + strconv.Itoa(user.Distance)
}

// There are two keys for each of this hash: farthest and nearest

// UpdateStatistics Attempts to update them both
func UpdateStatistics(database db.Interface, user model.User) {
	go updateStat(database, "farthest", statisticshash, user, farthestStrategy)
	go updateStat(database, "nearest", statisticshash, user, nearestStrategy)
}

// RetrieveDistance will rebuild the info from the stored {ISO code}-{distance} for the requested key
func RetrieveDistance(database db.Interface, key string) model.Statistic {
	str := database.Retrieve(statisticshash, key)
	info := strings.Split(str, "-")
	dst, err := strconv.Atoi(info[1])
	if err != nil {
		log.Printf("ERROR: cannot convert %s to int", info[1])
	}
	return model.Statistic{Country: info[0], Distance: dst}
}

func farthestStrategy(current, actual int) bool {
	return current < actual
}
func nearestStrategy(current, actual int) bool {
	return current > actual
}

func updateStat(database db.Interface, stat, hash string, user model.User, strategy func(int, int) bool) {
	if !database.Exists(hash, stat) {
		log.Printf("Storing first %s distance: %s with %d km from AR", stat, user.ISOCountry, user.Distance)
		database.Store(hash, stat, toValue(user))
		return
	}
	res := database.Retrieve(hash, stat)
	val := strings.Split(res, "-")

	currentStatDistance, err := strconv.Atoi(val[1])
	if err != nil {
		log.Printf("ERROR: cannot convert %s to int", val)
	}

	// If nearest is 0 we're done
	if stat == "nearest" && currentStatDistance == 0 {
		return
	}

	if strategy(currentStatDistance, user.Distance) {
		log.Printf("Storing new %s distance: %s with %d km from AR", stat, user.ISOCountry, user.Distance)
		database.Store(hash, stat, toValue(user))

	} else if currentStatDistance == user.Distance {
		// When equals take the one with (country) highest score

		var wg sync.WaitGroup
		wg.Add(2)

		var userIPCountryScore int
		go func() {
			defer wg.Done()
			userIPCountryScore = CountryBestScore(database, user.ISOCountry)
		}()

		var currentStatCountryScore int
		go func() {
			defer wg.Done()
			currentStatCountryScore = CountryBestScore(database, val[0])
		}()

		wg.Wait()

		if userIPCountryScore > currentStatCountryScore {
			log.Printf("Storing new %s country: %s because greater score. %s has a score of: %d while %s had %d", stat, user.ISOCountry, user.ISOCountry, userIPCountryScore, val[0], currentStatCountryScore)
			database.Store(hash, stat, toValue(user))
		}
	}
}
