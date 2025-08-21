package main

import (
	"database/sql"
	"log"
	"subscription-management/pkg/handler"

	"subscription-management/pkg/repository"

	"github.com/pressly/goose"
)

// MigrateDBRepo - executes the database migration file
func MigrateDBRepo(db *sql.DB) {

	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	if err := goose.Up(db, "/go/src/migrations"); err != nil {
		panic(err)
	}
}

func main() {

	var processConfig handler.ProcessConfig
	db, err := repository.SetUpDB()
	if err != nil {
		log.Fatal("failed to setup setup Database %w", err)
	}

	CurdRepo := repository.NewCurdRepo(db)
	processConfig.CurdRepo = CurdRepo
	MigrateDBRepo(db)

}
