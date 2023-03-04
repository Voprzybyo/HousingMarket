package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host     = "YOUR_HOST_NAME"
	port     = 5432
	user     = "YOUR_USER"
	password = "YOUR_PWD"
	dbname   = "YOUR_DATABASE_NAME"
)

func AddToDb(flatData []FlatData) {

	InfoLogger.Println("Writing data to database...")
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)

	defer db.Close()

	err = db.Ping()
	CheckError(err)

	for _, v := range flatData {
		insertDynStmt := `insert into "Flats"("Price", "Area", "Place", "PublicationDate", "FetchDate", "FetchHour", "InflationRate") values($1, $2, $3, $4, $5, $6, $7 )`
		_, e := db.Exec(insertDynStmt, v.Price, v.Area, v.Place, v.PublicationDate, v.FetchDate, v.FetchHour, v.InflationRate)
		CheckError(e)
	}

	InfoLogger.Println("Writing data to database finished!")
}
