package driver

import (
	"database/sql"
	"log"
	"os"

	"github.com/lib/pq"
)


var db *sql.DB


func handleError(err error) {
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func ConnectDB() (*sql.DB) {
	pgUrl, err := pq.ParseURL(os.Getenv("DB_URL"))
	handleError(err)

	// open pg connection
	db, err = sql.Open("postgres", pgUrl)
	handleError(err)

	// test to confirm that connection to the db works
	err = db.Ping()
	handleError(err)

	return db
}