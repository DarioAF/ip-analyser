package main

import (
	"log"
)

//RestCountriesResponse is the response from restCountries
type RestCountriesResponse struct {
	Alpha3Code string
	Latlng     [2]float64
}

func findLatlng(countries []RestCountriesResponse, countryCode3 string) [2]float64 {
	for _, n := range countries {
		if countryCode3 == n.Alpha3Code {
			return n.Latlng
		}
	}
	log.Panicf("Cant find [lat, lng] for iso country: %s", countryCode3)
	return [2]float64{0, 0}
}

func resolveCountryLocation(country IP2countryResponse) [2]float64 {
	res := []RestCountriesResponse{}
	uri := "https://restcountries.eu/rest/v2/name/" + country.CountryCode
	parseResponse(uri, &res)

	return findLatlng(res, country.CountryCode3)
}
