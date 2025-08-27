package config

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() error {
	var err error
	dbURL, err := Getenv("DB_URL")
	if err != nil {
		return err
	}
	DB, err = gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		return err
	}
	log.Println("Database connection established successfully")
	return nil
}
