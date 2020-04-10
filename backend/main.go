package main

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"

	"github.com/rjmarques/something-of-the-day/datastore"
	"github.com/rjmarques/something-of-the-day/service"
)

func main() {
	db := getDBConn()
	dao := datastore.NewDAO(db)
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")

	sotd := service.NewSomethingOfTheDay(clientID, clientSecret, dao)
	sotd.Start()
}

func getDBConn() *sql.DB {
	url := os.Getenv("POSTGRES_URL")
	db, err := sql.Open("postgres", url)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}
