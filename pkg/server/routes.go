package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/DarioAF/ip-analyser/pkg/db"
	"github.com/DarioAF/ip-analyser/pkg/external"
	"github.com/DarioAF/ip-analyser/pkg/model"
	"github.com/DarioAF/ip-analyser/pkg/service"
)

// ExternalUser is the user's input info
type ExternalUser struct {
	IP string
}

func isValidIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

func userHandler(w http.ResponseWriter, req *http.Request, database db.DBInterface) {
	var user ExternalUser

	if err := json.NewDecoder(req.Body).Decode(&user); err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, fmt.Sprintf(`{"error":"invalid input: %v"}`, err))
		return
	}

	if !isValidIP(user.IP) {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"error":"invalid ip"}`)
		return
	}

	country, err := external.ResolveCountry(user.IP)

	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()))
		return
	}

	start := time.Now()
	currentTime := start.Format("02/01/2006 15:04:05")

	var wg sync.WaitGroup
	wg.Add(2)

	var distance int
	go func() {
		defer wg.Done()
		distance = service.ResolveDistance(database, country)
	}()

	var isAWS bool
	go func() {
		defer wg.Done()
		isAWS = service.IsFromAWS(user.IP)
	}()

	wg.Wait()

	enhancedUser := model.User{
		IP:         user.IP,
		Time:       currentTime,
		Country:    country.CountryName,
		ISOCountry: country.CountryCode,
		Distance:   distance,
		IsAWS:      isAWS,
	}

	elapsed := time.Since(start)
	log.Printf("analysed ip: %s (%s) in %s", enhancedUser.IP, enhancedUser.Country, elapsed)

	go service.UpdateTrend(database, enhancedUser)
	go service.UpdateStatistics(database, enhancedUser)

	res, err := json.Marshal(enhancedUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, fmt.Sprintf("There was an error marshaling our user: %v", err))
	}

	w.Header().Add("Content-Type", "application/json")
	io.WriteString(w, string(res))
}

func distanceHandler(w http.ResponseWriter, r *http.Request, database db.DBInterface, impl string) {
	stat := service.RetrieveDistance(database, impl)

	res, err := json.Marshal(stat)
	if err != nil {
		log.Printf("There was an error marshaling our user! %err", err)
	}

	w.Header().Add("Content-Type", "application/json")
	io.WriteString(w, string(res))
}

func countryRequestsHandler(w http.ResponseWriter, r *http.Request, database db.DBInterface) {
	countryIso := r.URL.Path[len("/avg-requests/"):]
	avg := strconv.Itoa(service.CountryAvgRequests(database, countryIso))
	w.Header().Add("Content-Type", "application/json")
	io.WriteString(w, `{"avg":`+avg+"}")
}

func healthHandler(w http.ResponseWriter, r *http.Request, database db.DBInterface) {
	health := "UP & Running"
	w.Header().Add("Content-Type", "application/json")

	if database.Ping() != "PONG" {
		health = "something is wrong with the db, please check your connection with it"
		w.WriteHeader(http.StatusInternalServerError)
	}
	io.WriteString(w, fmt.Sprintf(`{"status":"%s"}`, health))
}