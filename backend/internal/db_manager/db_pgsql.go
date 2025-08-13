package db_manager

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"

	_ "github.com/lib/pq"
)

// wrapper to customise db query methods below.
// If we need to support different database types, change this to interface and make dependency injection as needed for different db types
type DB struct {
	db *sql.DB
}

var once sync.Once
var db *DB

func InitPgsqlConnection() *DB {
	once.Do(func() {
		//for now, using environment variables with fallback to defaults
		dbHost := getEnvOrDefault("DB_HOST", "localhost")
		dbPort := getEnvOrDefault("DB_PORT", "5432")
		dbUser := getEnvOrDefault("DB_USER", "user")
		dbPassword := getEnvOrDefault("DB_PASSWORD", "password")
		dbName := getEnvOrDefault("DB_NAME", "parser_db")

		dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			dbHost, dbPort, dbUser, dbPassword, dbName)

		log.Printf("Attempting to connect to database with host=%s port=%s user=%s dbname=%s", dbHost, dbPort, dbUser, dbName)

		var err error
		dbpgsql, err := sql.Open("postgres", dsn)
		if err != nil {
			log.Fatalf("Error opening database: %v", err)
		}

		if err = dbpgsql.Ping(); err != nil {
			log.Fatalf("Error connecting to database: %v", err)
		}

		log.Println("Successfully connected to PostgreSQL database")
		db = &DB{db: dbpgsql}
	})
	return db
}

func (d *DB) CreateRecord(ctx context.Context, query string, args ...interface{}) *sql.Row {
	tx := GetTransactionFromContext(ctx)
	if tx != nil {
		return tx.QueryRowContext(ctx, query, args...)
	}
	return d.db.QueryRowContext(ctx, query, args...)
}

func (d *DB) GetRecord(ctx context.Context, query string, args ...interface{}) *sql.Row {
	tx := GetTransactionFromContext(ctx)
	if tx != nil {
		return tx.QueryRowContext(ctx, query, args...)
	}
	return d.db.QueryRowContext(ctx, query, args...)
}

func (d *DB) UpdateRecord(ctx context.Context, query string, args ...interface{}) *sql.Row {
	tx := GetTransactionFromContext(ctx)
	if tx != nil {
		return tx.QueryRowContext(ctx, query, args...)
	}
	return d.db.QueryRowContext(ctx, query, args...)
}

func (d *DB) DeleteRecord(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	tx := GetTransactionFromContext(ctx)
	if tx != nil {
		return tx.ExecContext(ctx, query, args...)
	}
	return d.db.ExecContext(ctx, query, args...)
}

// GetRecords executes a query that returns multiple rows
func (d *DB) GetRecords(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	tx := GetTransactionFromContext(ctx)
	if tx != nil {
		return tx.QueryContext(ctx, query, args...)
	}
	return d.db.QueryContext(ctx, query, args...)
}

func (d *DB) Begin() (*sql.Tx, error) {
	return d.db.Begin()
}

// getEnvOrDefault returns the environment variable value or the default if not set
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func CloseDB() {
	if db != nil {
		err := db.db.Close()
		if err != nil {
			log.Printf("Error closing database: %v", err)
		} else {
			log.Println("Database connection closed")
		}
	}
}
