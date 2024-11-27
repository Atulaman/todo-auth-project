package main

import (
	"fmt"
	"log"
	"net/http"
	"todo-auth/database"
	"todo-auth/routes"

	_ "github.com/lib/pq"
)

func main() {
	db := database.Connect()
	defer db.Close()
	r := routes.Routes()
	fmt.Println("Server is running on http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", r))
}
