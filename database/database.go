package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"os"
)

var DBClient *sqlx.DB

func InitDB() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Print(err)
	}
	db, err := sqlx.Connect("postgres", os.Getenv("DB_KEY"))
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	DBClient = db
}
