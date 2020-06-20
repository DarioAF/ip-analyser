package main

import (
	"log"
)

//RestCountriesResponse is the response from restCountries
type RestCountriesResponse struct {
	Name       string
	Alpha3Code string
	Latlng     [2]float64
}

func findLatlng(countries []RestCountriesResponse, country IP2countryResponse) [2]float64 {
	for _, n := range countries {
		if country.CountryCode3 == n.Alpha3Code || country.CountryName == n.Name {
			return n.Latlng
		}
	}
	log.Printf("ERROR: cannot find [lat, lng] for iso country: %s or name: %s", country.CountryCode3, country.CountryName)
	return [2]float64{0, 0}
}

func resolveCountryLocation(country IP2countryResponse) [2]float64 {
	res := []RestCountriesResponse{}
	uri := "https://restcountries.eu/rest/v2/name/" + country.CountryCode
	parseResponse(uri, &res)

	return findLatlng(res, country)
}
