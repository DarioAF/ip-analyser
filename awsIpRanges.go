package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// IPRanges are the ip-ranges from AWS
type IPRanges struct {
	Prefixes      []IPv4Prefix
	Ipv6_prefixes []IPv6Prefix
}

// IPv4Prefix prefix in IPv4
type IPv4Prefix struct {
	Ip_prefix string
}

// IPv6Prefix prefix in IPv6
type IPv6Prefix struct {
	Ipv6_prefix string
}

func resolveAWSPrefixes() IPRanges {
	request := "https://ip-ranges.amazonaws.com/ip-ranges.json"
	prefixes := IPRanges{}

	webClient := &http.Client{Timeout: 10 * time.Second}
	res, err := webClient.Get(request)
	if err != nil {
		log.Printf("ERROR: there was an error getting: %s", request)
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&prefixes)
	if err != nil {
		log.Printf("ERROR: there was an error parsing the response from %s, status code: %d", request, res.StatusCode)
	}

	return prefixes
}
