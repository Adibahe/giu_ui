package main

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

func getName(id string) string {

	value, exists := recent[id]
	if exists {
		log.Println("Recieved Name via map: ", value.Name)
		return value.Name
	}

	querry := "SELECT name FROM FunctionSignatures WHERE rowid = ?"

	var name string
	err := db.QueryRow(querry, id).Scan(&name)
	if err != nil {
		log.Println("Failed to get Name for", id, "At getName in querry.go")
	}

	var msg message

	msg.Id = id
	msg.Name = name
	msg.Content = ""
	recent[id] = msg

	log.Println("Recieved Name: ", name)
	return name

}

func getDesc(id string, details *string) {

	value, _ := recent[id]

	if value.Content != "" {
		log.Println("Sent desc through map for: ", value.Name)
		*details = value.Content
		return
	}

	querry := `
		SELECT rd.document
		FROM FunctionSignatures fs
		JOIN RawDocument rd ON rd.url = fs.url
		WHERE fs.rowid = ?
	`
	err := db.QueryRow(querry, id).Scan(&value.Content)
	if err != nil {
		log.Println("Failed to get Description for", value.Name)
	}

	value.Content = htmlToReadableText(value.Content)

	log.Println("sent desc for: ", value.Name)
	*details = value.Content
	recent[id] = value
}

func connectDb() *sql.DB {

	db, err := sql.Open("sqlite", "ntdocs.sqlite3")
	if err != nil {
		log.Println("Failed to connect with DB")
	} else {
		log.Println("Connected to DB")
	}
	return db
}
