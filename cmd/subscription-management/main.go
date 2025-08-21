package main

import (
	"log"
	"subscription-management/pkg/handler"

	"subscription-management/pkg/repository"
)

func main() {

	var processConfig handler.ProcessConfig
	db, err := repository.SetUpDB()
	if err != nil {
		log.Fatal("failed to setup setup Database %w", err)
	}

	CurdRepo := repository.NewCurdRepo(db)
	processConfig.CurdRepo = CurdRepo
	MigrateDBRepo(db)
	initWebServer(processConfig)

}
