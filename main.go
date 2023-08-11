package main

import (
	"log"

	"github.com/EmilyOng/cvwo/backend/db"
	"github.com/EmilyOng/cvwo/backend/router"
	"github.com/joho/godotenv"
)

func main() {
	// Load any environment variables
	godotenv.Load()

	// Database setup
	err := db.Setup()
	if err != nil {
		log.Fatalln("Unable to setup database", err)
	}

	// Router setup
	router := router.Setup()
	router.Run()
}
