package main

import (
	"log"
	"net/http"
	"os"
)

func initServer() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/user", userHandler)
	http.HandleFunc("/nearest", nearestHandler)
	http.HandleFunc("/farthest", farthestHandler)
	http.HandleFunc("/avg-requests/", countryRequestsHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("** Service Started on Port " + port + " **")
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
