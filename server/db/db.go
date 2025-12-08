package db

import (
	"database/sql"
	"fmt"

	"github.com/InspectorGadget/realtime-polling-system/config"
	_ "github.com/lib/pq"
)

var db *sql.DB

func Connect() (*sql.DB, error) {
	connectionString := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		config.GetConfig("DB_USER"),
		config.GetConfig("DB_PASSWORD"),
		config.GetConfig("DB_HOST"),
		config.GetConfig("DB_PORT"),
		config.GetConfig("DB_NAME"),
	)

	var err error
	db, err = sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func GetDB() *sql.DB {
	return db
}
