package main

import (
	"RWTAPI/controller"
	"RWTAPI/sqldb"
	"log"
)

func main() {
	// Initialize the database connection
	sqldb.InitDB()

	log.Println("Server is running on port 8086")
	controller.InitHandlers()

}
