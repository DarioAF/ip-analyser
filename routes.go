package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"
)

// ExternalUser is the user's input info
type ExternalUser struct {
	IP string
}

// User is our enhanced user, containing all the aditional info
type User struct {
	IP         string
	Time       string
	Country    string
	ISOCountry string
	Distance   int
	IsAWS      bool
}

func deserializeUser(body io.ReadCloser) ExternalUser {
	decoder := json.NewDecoder(body)
	var user ExternalUser
	err := decoder.Decode(&user)
	if err != nil {
		panic(err)
	}
	return user
}

func isValidIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

func userHandler(w http.ResponseWriter, req *http.Request) {

	user := deserializeUser(req.Body)
	if isValidIP(user.IP) {
		country, err := resolveCountry(user.IP)

		if err == nil {
			start := time.Now() // measure start
			currentTime := start.Format("02/01/2006 15:04:05")
			distanceChan := make(chan int)
			isAwsChan := make(chan bool)

			go resolveDistance(country, distanceChan)
			go isFromAWS(user.IP, isAwsChan)

			distance := <-distanceChan
			isAWS := <-isAwsChan

			enhancedUser := User{
				user.IP,
				currentTime,
				country.CountryName,
				country.CountryCode,
				distance,
				isAWS}

			elapsed := time.Since(start) //measure stop
			log.Printf("analysed ip: %s (%s) in %s", enhancedUser.IP, enhancedUser.Country, elapsed)

			go updateTrend(enhancedUser)
			go updateStatistics(enhancedUser)

			res, err := json.Marshal(enhancedUser)
			if err != nil {
				log.Panicf("There was an error marshaling our user! ", err)
			}

			w.Header().Add("Content-Type", "application/json")
			io.WriteString(w, string(res))

		} else {
			w.Header().Add("Content-Type", "application/json")
			io.WriteString(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()))
		}
	} else {
		w.Header().Add("Content-Type", "application/json")
		io.WriteString(w, `{"error":"invalid ip"}`)
	}
}

func nearestHandler(w http.ResponseWriter, r *http.Request) {
	nearest := retrieveNearest()
	res, err := json.Marshal(nearest)
	if err != nil {
		log.Panicf("There was an error marshaling our user! ", err)
	}

	w.Header().Add("Content-Type", "application/json")
	io.WriteString(w, string(res))
}

func farthestHandler(w http.ResponseWriter, r *http.Request) {
	farthest := retrieveFarthest()
	res, err := json.Marshal(farthest)
	if err != nil {
		log.Panicf("There was an error marshaling our user! ", err)
	}

	w.Header().Add("Content-Type", "application/json")
	io.WriteString(w, string(res))
}

func countryRequestsHandler(w http.ResponseWriter, r *http.Request) {
	countryIso := r.URL.Path[len("/avg-requests/"):]
	avg := strconv.Itoa(countryAvgRequests(countryIso))
	w.Header().Add("Content-Type", "application/json")
	io.WriteString(w, `{"avg":`+avg+"}")
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	io.WriteString(w, `{"status":"ok"}`)
}
