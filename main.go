package main

import (
	"log"

	"github.com/EmilyOng/tusk-manager/backend/db"
	"github.com/EmilyOng/tusk-manager/backend/router"
	"github.com/joho/godotenv"
)

func main() {
	// Load any environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Unable to load environment variables", err)
	}

	// Database setup
	err = db.Setup()
	if err != nil {
		log.Fatalln("Unable to setup database", err)
	}

	// Router setup
	router := router.Setup()
	err = router.Run()
	if err != nil {
		log.Fatalln("Unable to run the router", err)
	}
}
