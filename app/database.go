package app

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

func NewDB() *sql.DB {
	log := logrus.Logger{}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PARAMS"),
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Panicf("failed to open db connection: %v", err)
	}

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(10)
	db.SetConnMaxIdleTime(10 * time.Minute)
	db.SetConnMaxLifetime(60 * time.Minute)
	return db
}
