package db

import (
	"database/sql"
	"log"
	"os"
	_ "modernc.org/sqlite"
)

const (
	schemaCreateTable = `CREATE TABLE IF NOT EXISTS scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date CHAR(8) NOT NULL DEFAULT '',
		title VARCHAR(256) NOT NULL DEFAULT '',
		comment TEXT NOT NULL DEFAULT '',
		repeat VARCHAR(128) NOT NULL DEFAULT ''
	)`

	schemaCreateIdx = `CREATE INDEX IF NOT EXISTS idx_date ON scheduler (date);`
)

var db *sql.DB

func Init() error {
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}
	log.Printf("Using database file: %s", dbFile)

	_, err := os.Stat(dbFile)

	var install bool

	if err != nil {
		install = true
	}

	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		log.Printf("error opening database %s: %s", dbFile, err)
		return err
	}

	if err = db.Ping(); err != nil {
		log.Printf("database connection error: %s", err)
		return err
	}

	if install {
		log.Printf("there is no database, create a new one: %s", dbFile)

		_, err := db.Exec(schemaCreateTable)
		if err != nil {
			log.Printf("error creating table: %s", err)
			return err
		}
		_, err = db.Exec(schemaCreateIdx)
		if err != nil {
			log.Printf("error creating index: %s", err)
			return err
		}

		log.Printf("database initialized successfully: %s", dbFile)
	}

	return nil
}

func Close() {
	if db != nil {
		err := db.Close()
		if err != nil {
			log.Printf("error closing database")
		} else {
			log.Panicln("database connection closed")
		}
	}
}
