package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq" // PostgreSQL driver
)

var TODO *sql.DB

// Initialize Database
func ConnectDB() *sql.DB {
	//db string
	connStr := "host=localhost port=5432 user=postgres password=rx dbname=todo-multi sslmode=disable"

	// Open a connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	// Ensure the connection is not nil
	if db == nil {
		log.Fatal("Database connection is nil")
	}

	fmt.Println("Connected to the database successfully!")

	err = migrateUp(db)
	if err != nil {
		log.Fatal("Migrations failed", err)
	}

	fmt.Println("Migrations completed successfully!")

	TODO = db
	// utils.SetDB(db)
	return db

}

func migrateUp(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://database/migrations",
		"postgres", driver)

	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
