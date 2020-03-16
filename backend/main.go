package main

import (
	"database/sql"
	"fmt"
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

const (
	port   = 5432
	user   = "postgres"
	dbname = "somethingoftheday"
)

func getDBConn() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		os.Getenv("HOST"), port, user, os.Getenv("DB_PASS"), dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}

