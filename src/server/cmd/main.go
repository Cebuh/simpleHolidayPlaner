package main

import (
	"database/sql"
	"log"

	"github.com/cebuh/simpleHolidayPlaner/cmd/api"
	"github.com/cebuh/simpleHolidayPlaner/config"
	"github.com/cebuh/simpleHolidayPlaner/db"
	"github.com/go-sql-driver/mysql"
)

func main() {
	db, err := db.NewMySqlStorage(mysql.Config{
		User:                 config.Envs.DBUser,
		Passwd:               config.Envs.DBPassword,
		Addr:                 config.Envs.DBAddress,
		DBName:               config.Envs.DBName,
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
