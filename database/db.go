package database

import (
	//"database/sql"
	"fmt"
	//"log"
	"os"

	"github.com/jmoiron/sqlx"

	log "todo-auth/logging"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

var TODO *sqlx.DB

func migrateUp(db *sqlx.DB) {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		log.Logging(err, "Failed to connect to the database", 500, "fatal", nil)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://database/migrations",
		"postgres", driver)
	if err != nil {
		log.Logging(err, "Failed to connect to the database", 500, "fatal", nil)
	}
	//err =m.Up() // or m.Steps(2) if you want to explicitly set the number of migrations to run
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Logging(err, "Failed to migrate the database", 500, "fatal", nil)
	}
}
func Connect() {
	//connStr := "host=localhost port=5433 user=postgres password=rx dbname=todo sslmode=disable"
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), "postgres", "rx", "todo")

	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		log.Logging(err, "Failed to connect to the database", 500, "fatal", nil)
	}

	err = db.Ping()
	if err != nil {
		log.Logging(err, "Failed to connect to the database", 500, "fatal", nil)
	}
	fmt.Println("Connected to the database successfully!")
	log.Logging(nil, "Connected to the database successfully!", 200, "info", nil)
	migrateUp(db)
	TODO = db
}
func ShutDownDb() error {
	return TODO.Close()
}
