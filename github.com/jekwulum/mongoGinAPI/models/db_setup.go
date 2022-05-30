package models

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
)

var DB *gorm.DB

func ConnectPostGRES_DB() {
	envErr := godotenv.Load(".env")

	if envErr != nil {
		log.Fatal("Error loading .env file")
	}

	dbDriver := os.Getenv("DB_DRIVER")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("PG_DB_PORT")
	db_name := os.Getenv("PG_DB_NAME")
	user := os.Getenv("PG_DB_USER")
	password := os.Getenv("PG_DB_PASSWORD")
	DB_URL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", host, port, user, db_name, password)
	
	database, dbErr := gorm.Open(dbDriver, DB_URL)

	if dbErr != nil {
		fmt.Println("dbErr------------->", dbErr)
		panic("Failed to connect to database!")
	}

	fmt.Println("Connected to Postgres Database!")
	database.AutoMigrate(&Cust_User{})

	DB = database
}