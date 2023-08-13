package db

import (
	"log"
	"os"

	"github.com/EmilyOng/cvwo/backend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Setup() (err error) {
	DB_URL := os.Getenv("DATABASE_URL")
	DB, err = gorm.Open(postgres.Open(DB_URL), &gorm.Config{})
	if err != nil {
		log.Fatalln("Invalid database configuration")
		return
	}

	err = DB.AutoMigrate(
		&models.User{},
		&models.Board{},
		&models.Task{},
		&models.Tag{},
		&models.State{},
		&models.Member{},
	)
	if err != nil {
		log.Fatalln("Unable to migrate database")
		return
	}

	return
}
