package main

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var db *sqlx.DB

func initDB() {
  var err error
  dsn := "postgres://erp_user:erp_password@localhost:5432/erp_asset?sslmode=disable"
  db, err = sqlx.Connect("postgres", dsn)
  if err != nil {
    log.Fatalln("DB connection error:", err)
  }
}
