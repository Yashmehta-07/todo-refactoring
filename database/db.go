package database

import (
	"fmt"
	"os"
	"todo/logging"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // PostgreSQL driver
)

var TODO *sqlx.DB

// Initialize Data1base
func ConnectDB() *sqlx.DB {
	//db string
	// connStr := "host=localhost port=5432 user=postgres password=rx dbname=todo-multi sslmode=disable"
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	// Open a connection
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		logging.Log(err, "Error connecting to the database", "fatal", 500, nil)
	}

	fmt.Println("Connected to the database successfully!")

	err = migrateUp(db)
	if err != nil {
		logging.Log(err, "Error running migrations", "fatal", 500, nil)
	}

	fmt.Println("Migrations completed successfully!")

	TODO = db
	// utils.SetDB(db)
	return db

}

func migrateUp(db *sqlx.DB) error {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		logging.Log(err, "Error creating migration driver", "fatal", 500, nil)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://database/migrations",
		"postgres", driver)

	if err != nil {
		logging.Log(err, "Error creating migration instance", "fatal", 500, nil)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
