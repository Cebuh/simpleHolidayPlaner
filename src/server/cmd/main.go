package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/cebuh/simpleHolidayPlaner/cmd/api"
	"github.com/cebuh/simpleHolidayPlaner/db"
	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Fehler beim Laden der .env Datei: %v", err)
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	address := fmt.Sprintf("%s:%s", dbHost, dbPort)

	db, err := db.NewMySqlStorage(mysql.Config{
		User:                 dbUser,
		Passwd:               dbPassword,
		Addr:                 address,
		DBName:               dbName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	})

	if err != nil {
		log.Fatal(err)
	}

	initStorage(db)
	server := api.NewServer(":8080", db)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

func initStorage(db *sql.DB) {
	log.Println("Check Database connection...")
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Database successfully connected!")
}
