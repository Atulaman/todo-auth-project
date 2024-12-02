package main

import (

	//"log"
	"net/http"
	"todo-auth/database"
	"todo-auth/routes"

	log "todo-auth/logging"

	_ "todo-auth/docs"

	_ "github.com/lib/pq"
)

// @title To-Do API
// @version 1.0
// @description This is a sample server for a To-Do app.
// @host localhost:8081
// @BasePath /
func main() {
	database.Connect()
	defer database.ShutDownDb()
	r := routes.Routes()
	log.Logger.Info("Server is running on http://localhost:8081")
	log.Logger.Fatal(http.ListenAndServe(":8081", r))
}
