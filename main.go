package main

import "github.com/JovanZdravkovic/TaskJournalBackend/api"

func main() {
	router := api.NewRouter(":8080")
	router.ConfigureRoutes()
	router.ListenAndServe()
}
