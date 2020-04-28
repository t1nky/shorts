package main

import (
	"fmt"
	"log"
	"os"

	"shorts/database"
	_ "shorts/docs"
	"shorts/models"
	"shorts/router"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
)

// InitDatabase : Initialize database
func InitDatabase() (*gorm.DB, error) {
	// Load env variables
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Turn on logging if needed by: database.DB.LogMode(true)

	// Connect to the DB
	db, err := gorm.Open("postgres", fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		user,
		password,
		host,
		port,
		dbName,
	))
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Shortlink{})
	db.AutoMigrate(&models.ShortlinkUse{})

	database.DB = db

	return db, nil
}

func main() {

	// Init local env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := InitDatabase()
	if err != nil {
		fmt.Println("Cannot connect to the database:" + err.Error())
		return
	}
	defer db.Close()

	// Initialize WebServer
	r := router.SetupRouter()

	r.Run()
}
