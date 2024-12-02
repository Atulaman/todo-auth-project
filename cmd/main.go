package main

import (

	//"log"
	"net/http"
	"todo-auth/database"
	"todo-auth/routes"

	log "todo-auth/logging"

	_ "github.com/lib/pq"
)

//	func init() {
//		log.Init()
//	}
func main() {
	database.Connect()
	defer database.ShutDownDb()
	r := routes.Routes()
	//fmt.Println("Server is running on http://localhost:8081")
	log.Logger.Info("Server is running on http://localhost:8081")
	log.Logger.Fatal(http.ListenAndServe(":8081", r))
}
