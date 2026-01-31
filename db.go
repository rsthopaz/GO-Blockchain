package main

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var db *sqlx.DB

func initDB() {
	var err error

	dsn := `
		host=localhost
		port=5432
		user=erp_user
		password=erp123
		dbname=erp_asset
		sslmode=disable
	`

	db, err = sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalln("DB connection error:", err)
	}
}
