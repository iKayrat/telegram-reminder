package controllers

import (
	"database/sql"
	"log"
)

func DBconnection(dbsource string) *sql.DB {
	// connStr := "user=pg dbname=pqgotest sslmode=verify-full"
	db, err := sql.Open("postgres", dbsource)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	return db
}
