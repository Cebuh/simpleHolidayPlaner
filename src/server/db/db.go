package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
)

func NewMySqlStorage(cfg mysql.Config) (*sql.DB, error) {
	cfg.MultiStatements = true
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?checkConnLiveness=false&multiStatements=true&parseTime=true&maxAllowedPacket=0", cfg.User, cfg.Passwd, cfg.Addr, cfg.DBName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	return db, nil
}
