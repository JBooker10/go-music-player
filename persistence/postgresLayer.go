package persistence


import (
	"database/sql"
	"log"
	"os"

	"github.com/lib/pq"
)

var sqlDB *sql.DB

func PostgreSQLConnection() *sql.DB {
	postgresURL, err := pq.ParseURL(os.Getenv("POSTGRES_URL"))
	if err != nil {
		log.Fatal(err)
	}

	sqlDB, err = sql.Open("postgres", postgresURL)
	if err != nil {
		log.Fatal(err)
	}

	err = sqlDB.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("PostgreSQL connected")
	return sqlDB
}

