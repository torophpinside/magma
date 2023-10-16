package service

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func NewSqliteDatabase(regenerate bool) *sql.DB {
	if regenerate {
		os.Remove("./dorks.db")
	}

	db, err := sql.Open("sqlite3", "./dorks.db")
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func RunSqliteMigration(conn *sql.DB) {
	queries := []string{
		"CREATE TABLE IF NOT EXISTS dorks (id INTEGER PRIMARY KEY, dork TEXT UNIQUE, score INTEGER)",
	}

	for _, query := range queries {
		_, err := conn.Exec(query)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func DBConn() *sql.DB {
	return db
}
