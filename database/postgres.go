package database

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func InitDbPostgres() (*sql.DB, error) {
	host := os.Getenv("POSTGRES_HOST")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DBNAME")

	port, _ := strconv.Atoi(os.Getenv("POSTGRES_PORT"))
	if port == 0 {
		port = 5432
	}

	conn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbName)

	slog.Debug(conn)

	db, err := sqlx.Open("postgres", conn)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	slog.Info("PostgreSQL successfully connected", "port", port)

	return db.DB, err
}
