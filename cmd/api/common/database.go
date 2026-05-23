package common

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type DBConfig struct {
	DB *pgxpool.Pool
}

func ConnectToDatabase() (db *DBConfig, error error) {
	// Implementation for connecting to the database
	err := godotenv.Load()

	if err != nil {
		panic("Error loading .env file: " + err.Error())
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	dsn := "postgres://" + user + ":" + password + "@" + host + ":" + port + "/" + dbname

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	return &DBConfig{DB: pool}, nil
}
