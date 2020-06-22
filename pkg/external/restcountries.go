package external

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
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

func ResolveCountryLocation(country IP2countryResponse) [2]float64 {
	locations := []RestCountriesResponse{}
	request := "https://restcountries.eu/rest/v2/name/" + country.CountryCode

	var webClient = &http.Client{Timeout: 10 * time.Second}
	res, err := webClient.Get(request)
	if err != nil {
		log.Printf("ERROR: there was an error getting: %s", request)
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&locations)
	if err != nil {
		log.Printf("ERROR: there was an error parsing the response from %s, status code: %d", request, res.StatusCode)
	}

	return findLatlng(locations, country)
}
