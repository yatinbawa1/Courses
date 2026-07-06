package database

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func nn(a error, b error) error {
	if a != nil {
		return a
	}

	return b
}

func ConnectDataBase() error {

	host := os.Getenv("DB_HOST")
	port, _ := strconv.Atoi(os.Getenv("DB_PORT"))
	user := os.Getenv("DB_USER")
	dbname := os.Getenv("DB_NAME")
	pass := os.Getenv("DB_PASSWORD")

	psqlConnectionString := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable", host, port, user, dbname, pass)

	db, errSql := sql.Open("postgres", psqlConnectionString)
	err := db.Ping()

	if errSql != nil || err != nil {
		return fmt.Errorf("Unable to Connect to Database because of the following error: \n\033[31m%w\033[0m", nn(errSql, err))
	} else {
		DB = db
		return nil
	}
}
