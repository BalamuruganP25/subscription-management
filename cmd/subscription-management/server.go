package main

import (
	"database/sql"
	"log"
	"net/http"

	"subscription-management/pkg/handler"
	"subscription-management/pkg/handler/customer"
	"subscription-management/pkg/handler/tax"
	"subscription-management/pkg/handler/user"
	"subscription-management/pkg/handler/webhook"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
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

func initWebServer(config handler.ProcessConfig) *http.Server {

	log.Println("🚀 Server is starting on :8089...")
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Route("/v1", func(v1 chi.Router) {
		// user part
		v1.Post("/api/users", user.CreateUser(&config))
		v1.Get("/api/users/{id}", user.GetUserById(&config))
		v1.Patch("/api/users/{id}", user.UpdateUserById(&config))
		v1.Delete("/api/users/{id}", user.DeleteUserById(&config))

		// customer part
		v1.Post("/api/customers", customer.CreateCustomer(&config))
		v1.Post("/api/subscriptions", customer.CreateSubscription(&config))

		// tax part
		v1.Get("/api/tax/{country}/{state}/{city}/{amount}/{postal_code}", tax.GetTax(&config))

		// webhook part
		v1.Post("/api/webhook", webhook.WebhookHandler(&config))

	})

	return &http.Server{
		Addr:    ":8089",
		Handler: r,
	}

}
