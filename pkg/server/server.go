package server

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/DarioAF/ip-analyser/pkg/db"
)

// Init initialize the server app
func Init() {
	var redis db.Interface = &db.RedisConnector

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		healthHandler(w, r, redis)
	})

	http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		userHandler(w, r, redis)
	})

	http.HandleFunc("/nearest", func(w http.ResponseWriter, r *http.Request) {
		distanceHandler(w, r, redis, "nearest")
	})

	http.HandleFunc("/farthest", func(w http.ResponseWriter, r *http.Request) {
		distanceHandler(w, r, redis, "farthest")
	})

	http.HandleFunc("/avg-requests/", func(w http.ResponseWriter, r *http.Request) {
		countryRequestsHandler(w, r, redis)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("** Service Started on Port " + port + " **")
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func serveResponse(w http.ResponseWriter, status int, body string) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	io.WriteString(w, body)
}
