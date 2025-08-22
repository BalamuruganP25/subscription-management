package main

import (
	"fmt"
	"log"
	"os"
	"subscription-management/pkg/handler"

	"subscription-management/pkg/repository"
)

func main() {

	var processConfig handler.ProcessConfig
	db, err := repository.SetUpDB()
	if err != nil {
		log.Fatal("failed to setup setup Database %w", err)
	}
	fmt.Println(os.Getenv("STRIPE_SECRET_KEY"))
	CurdRepo := repository.NewCurdRepo(db)
	processConfig.CurdRepo = CurdRepo
	MigrateDBRepo(db)
	initWebServer(processConfig)

}
