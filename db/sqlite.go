package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func init() {
	log.Print("Creating database... ")

	file, err := os.Create("db.sqlite")

	if err != nil {
		log.Println("FAILED")
		log.Fatal(err.Error())
	}

	file.Close()
	log.Println("Done")

	db, _ = sql.Open("sqlite3", "./db.sqlite")
	// defer db.Close()

	// create tables
	createSql := `--sql
		CREATE TABLE IF NOT EXISTS friends (
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			link TEXT NOT NULL UNIQUE,
			mutual INTEGER,
			unfriended BOOLEAN DEFAULT FALSE
		) `

	log.Print("Creating tables... ")
	query, err := db.Prepare(createSql)

	if err != nil {
		log.Println("Failed")
		log.Fatalln(err.Error())
	}

	query.Exec()
	log.Println("Done")
}

func InsertToFriend(name string, mutual int, link string) {
	insertSql := "INSERT INTO friends (name, link, mutual) VALUES (?, ?, ?) "
	query, err := db.Prepare(insertSql)
	if err != nil {
		log.Fatalln(err.Error())
	}

	_, err = query.Exec(name, link, mutual)
	if err != nil {
		if err.Error() == "UNIQUE constraint failed: friends.link" {
			log.Println(link, "found duplicate")
			return
		}

		log.Fatalln(err.Error())
	}
	// log.Println("Inserted friend ", name)
}
