package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	_ "github.com/lib/pq" // pq driver
)

var (
	db   *sql.DB
	once sync.Once
)

// GetDB returns a singleton *sql.DB instance
func GetDB() *sql.DB {
	once.Do(func() {
		dbUser := os.Getenv("POSTGRES_USER")
		dbPassword := os.Getenv("POSTGRES_PASSWORD")
		dbHost := os.Getenv("POSTGRES_HOST")
		dbPort := os.Getenv("POSTGRES_PORT")
		dbName := os.Getenv("POSTGRES_DB")

		if dbUser == "" || dbPassword == "" || dbHost == "" || dbPort == "" || dbName == "" {
			log.Fatal("Postgres environment variables not set")
		}

		// PQ connection string format
		dsn := fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable",
			dbUser, dbPassword, dbHost, dbPort, dbName,
		)

		var err error
		db, err = sql.Open("postgres", dsn)
		if err != nil {
			log.Fatalf("failed to open db: %v", err)
		}

		// Set connection pool
		db.SetMaxOpenConns(25)
		db.SetMaxIdleConns(25)
		db.SetConnMaxLifetime(5 * time.Minute)

		// Ping DB to verify connection
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := db.PingContext(ctx); err != nil {
			log.Fatalf("failed to ping db: %v", err)
		}

		log.Println("Connected to Postgres successfully")
	})

	return db
}

// CloseDB closes the database connection
func CloseDB() error {
	if db != nil {
		return db.Close()
	}
	return nil
}
