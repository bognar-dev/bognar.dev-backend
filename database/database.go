package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	supa "github.com/nedpals/supabase-go"
	"os"
)

var DBClient *sqlx.DB
var SBClient *supa.Client

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

func InitSupabase() {
	supabaseUrl := os.Getenv("SUPA_URL")
	supabaseKey := os.Getenv("SUPA_KEY")
	Supa := supa.CreateClient(supabaseUrl, supabaseKey)
	SBClient = Supa

}
