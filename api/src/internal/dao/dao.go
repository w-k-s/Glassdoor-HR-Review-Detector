package dao

import (
	"database/sql"
	"log"
)

func Must(db *sql.DB, err error) *sql.DB {
	if err != nil {
		log.Panic("failed to open sqlite db")
	}
	return db
}
