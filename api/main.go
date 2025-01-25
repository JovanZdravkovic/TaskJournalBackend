package api

func main() {
	router := NewRouter(":8080")
	router.ConfigureRoutes()
	router.ListenAndServe()
}
