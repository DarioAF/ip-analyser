package main

import (
	"fmt"
)

// IP2countryResponse is the response from the ip2country service
type IP2countryResponse struct {
	CountryCode  string
	CountryCode3 string
	CountryName  string
}

func resolveCountry(ip string) (IP2countryResponse, error) {
	country := IP2countryResponse{}
	request := "https://api.ip2country.info/ip?" + ip
	parseResponse(request, &country)

	if country.CountryCode3 == "" {
		return country, fmt.Errorf("Couldn't find any country for ip: $s", ip)
	}

	return country, nil
}
