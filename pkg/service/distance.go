package service

import (
	"log"
	"math"
	"strconv"

	"github.com/DarioAF/ip-analyser/pkg/db"
	"github.com/DarioAF/ip-analyser/pkg/external"
)

// Distances will be stored with hash: distance-{ISO}
// Since its allways from AR, is static
var distancehash string = "distance-AR"

// ResolveDistance will calc the distance only if it´s the first time the pair is seen, note that AR-AR is allways 0
func ResolveDistance(database db.Interface, country external.IP2countryResponse) int {
	if country.CountryCode == "AR" {
		return 0
	}

	if database.Exists(distancehash, country.CountryCode) {
		str := database.Retrieve(distancehash, country.CountryCode)
		res, err := strconv.Atoi(str)
		if err != nil {
			log.Printf("Cannot parse %s to int", str)
		}
		return res
	}

	location := external.ResolveCountryLocation(country)
	distance := int(distanceFromARGinKM(location[0], location[1]))
	database.Store(distancehash, country.CountryCode, strconv.Itoa(distance))

	return distance
}

func distanceFromARGinKM(lat float64, lng float64) float64 {
	// lat: -34.0, lng: -64 are from Argentina
	return calcDistance(-34.0, -64.0, lat, lng, "K")
}

//:::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
//:::  This routine calculates the distance between two points (given the     :::
//:::  latitude/longitude of those points). It is being used to calculate     :::
//:::  the distance between two locations using GeoDataSource (TM) prodducts  :::
//:::                                                                         :::
//:::  Definitions:                                                           :::
//:::    South latitudes are negative, east longitudes are positive           :::
//:::                                                                         :::
//:::  Passed to function:                                                    :::
//:::    lat1, lon1 = Latitude and Longitude of point 1 (in decimal degrees)  :::
//:::    lat2, lon2 = Latitude and Longitude of point 2 (in decimal degrees)  :::
//:::    unit = the unit you desire for results                               :::
//:::           where: 'M' is statute miles (default)                         :::
//:::                  'K' is kilometers                                      :::
//:::                  'N' is nautical miles                                  :::
//:::                                                                         :::
//::: thanks to: GeoDataSource.com                                            :::
//::: source: https://www.geodatasource.com/developers/go                     :::
//:::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
func calcDistance(lat1 float64, lng1 float64, lat2 float64, lng2 float64, unit ...string) float64 {
	const PI float64 = 3.141592653589793

	radlat1 := float64(PI * lat1 / 180)
	radlat2 := float64(PI * lat2 / 180)

	theta := float64(lng1 - lng2)
	radtheta := float64(PI * theta / 180)

	dist := math.Sin(radlat1)*math.Sin(radlat2) + math.Cos(radlat1)*math.Cos(radlat2)*math.Cos(radtheta)

	if dist > 1 {
		dist = 1
	}

	dist = math.Acos(dist)
	dist = dist * 180 / PI
	dist = dist * 60 * 1.1515

	if len(unit) > 0 {
		if unit[0] == "K" {
			dist = dist * 1.609344
		} else if unit[0] == "N" {
			dist = dist * 0.8684
		}
	}

	return dist
}
