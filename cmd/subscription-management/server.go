package main

import (
	"database/sql"
	"log"
	"net/http"

	"subscription-management/pkg/handler"
	"subscription-management/pkg/handler/user"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
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

func initWebServer(config handler.ProcessConfig) {

	log.Println("🚀 Server is starting on :8089...")
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Route("/v1", func(v1 chi.Router) {
		v1.Post("/user", user.CreateUser(&config))
		v1.Get("/user/{userID}", user.GetUserById(&config))
		v1.Patch("/user/{id}", user.UpdateUserById(&config))
		v1.Delete("/user/{id}", user.DeleteUserById(&config))
	})

	// ✅ Check for ListenAndServe error
	if err := http.ListenAndServe(":8089", r); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
