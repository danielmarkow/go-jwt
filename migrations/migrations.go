package migrations

import (
	"database/sql"
	"log"
)

var creates = map[string]string{
	"createUserTbl": `CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			email TEXT UNIQUE NOT NULL,
            password TEXT NOT NULL
		);`,
}

func CreateTables(db *sql.DB) {
	for stmt, create := range creates {
		_, err := db.Exec(create)
		if err != nil {
			log.Printf("error executing create table statement %s: %s \n", stmt, err.Error())
		} else {
			log.Printf("successfully ran %s \n", stmt)
		}
	}
}
