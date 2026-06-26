package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/jules/ecommerce-cqrs-go/internal/api"
	"github.com/jules/ecommerce-cqrs-go/internal/eventstore"
	"github.com/jules/ecommerce-cqrs-go/internal/projection"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/ecommerce?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Wait for DB
	for i := 0; i < 20; i++ {
		err = db.Ping()
		if err == nil {
			break
		}
		log.Printf("Waiting for database... %v", err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatal("Could not connect to database")
	}

	if err := applySchemas(db); err != nil {
		log.Fatal(err)
	}

	store := eventstore.NewStore(db)
	projectionist := projection.NewProjectionist(db)
	commandHandler := api.NewCommandHandler(store, projectionist)
	queryHandler := api.NewQueryHandler(db)

	r := mux.NewRouter()
	apiRouter := api.NewRouter(commandHandler, queryHandler)
	apiRouter.RegisterRoutes(r)

	fmt.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func applySchemas(db *sql.DB) error {
	schemas := []string{
		"internal/eventstore/schema.sql",
		"internal/projection/read_schema.sql",
	}

	for _, schemaFile := range schemas {
		content, err := os.ReadFile(schemaFile)
		if err != nil {
			return err
		}
		if _, err := db.Exec(string(content)); err != nil {
			return err
		}
	}
	return nil
}
