package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// IP2countryResponse is the response from the ip2country service
type IP2countryResponse struct {
	CountryCode  string
	CountryCode3 string
	CountryName  string
}

// These represent: Europe, Asia, North America, Africa, Oceania, Antarctica, South America
var continents = [7]string{"EU", "AS", "NA", "AF", "OC", "AN", "SA"}

func (cr IP2countryResponse) isContinent() bool {
	for _, c := range continents {
		if cr.CountryCode == c {
			return true
		}
	}
	return false
}

func resolveCountry(ip string) (IP2countryResponse, error) {
	country := IP2countryResponse{}
	request := "https://api.ip2country.info/ip?" + ip

	var webClient = &http.Client{Timeout: 10 * time.Second}
	res, err := webClient.Get(request)
	if err != nil {
		log.Printf("ERROR: there was an error getting: %s", request)
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&country)
	if err != nil {
		log.Printf("ERROR: there was an error parsing the response from %s, status code: %d", request, res.StatusCode)
	}

	if country.CountryCode3 == "" {
		return country, fmt.Errorf("couldn't find any country for ip: %s", ip)
	}

	if country.isContinent() {
		return country, fmt.Errorf("%s is a acontinent reather than a country, we can't handle that", country.CountryName)
	}

	return country, nil
}
