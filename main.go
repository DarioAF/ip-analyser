package main

import (
	"log"
)

// docker-compose up --build
func main() {
	log.Println("Hello World")

	pong := db.Ping()
	log.Printf("Executing ping command to redis... %s!", pong)

	initServer()
}
