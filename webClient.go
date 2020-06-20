package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

var webClient = &http.Client{Timeout: 10 * time.Second}

func parseResponse(url string, target interface{}) {
	resp, urlErr := webClient.Get(url)
	if urlErr != nil {
		log.Panicf("There was an error getting the url: %s", url)
	}
	defer resp.Body.Close()

	parErr := json.NewDecoder(resp.Body).Decode(target)
	if parErr != nil {
		log.Panicf("There was an error parsing the response from %s, status code: %d", url, resp.StatusCode)
	}
}
