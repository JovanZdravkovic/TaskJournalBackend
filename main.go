package main

import (
	"context"
	"log"
	"os"

	"github.com/JovanZdravkovic/TaskJournalBackend/api"
	"github.com/JovanZdravkovic/TaskJournalBackend/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	dbConnectionString := os.Getenv("DATABASE_URL")
	dbPool, err := pgxpool.New(context.Background(), dbConnectionString)
	if err != nil {
		panic("Could not establish a database connection pool")
	} else {
		log.Println("Successfully connected to database")
	}
	defer dbPool.Close()
	dbService := db.NewDatabaseService(dbPool)

	router := api.NewRouter(":8080")
	router.ConfigureRoutes(dbService)
	router.ListenAndServe()
}
