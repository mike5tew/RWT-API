package main

import (
	"RWTAPI/controller"
	"RWTAPI/sqldb"
	"log"
	"net/http"
)

func main() {
	// Initialize the database connection
	sqldb.InitDB("user:password@tcp(localhost:3306)/dbname")

	// Initialize handlers
	controller.InitHandlers()

	log.Println("Server is running on port 8086")
	log.Fatal(http.ListenAndServe(":8086", nil))
}
