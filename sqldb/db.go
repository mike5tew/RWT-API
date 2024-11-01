package sqldb

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type Connection struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// DB is a global variable to hold db connection
var DB *sql.DB

func ConnToString() string {
	envMap, err := godotenv.Read("config/.env")
	if err != nil {
		fmt.Printf("Error loading .env file")
		os.Exit(1)
		return ""
	}

	conStr := fmt.Sprintf("%s:%s@tcp(localhost:%s)/%s",
		envMap["MYSQL_USER"], envMap["MYSQL_ROOT_PASSWORD"], envMap["MYSQL_PORT"], envMap["MYSQL_DATABASE"])
	fmt.Println("HERE WE ARE " + conStr)
	return conStr

}

// ConnectDB opens a connection to the database
func ConnectDB() {

	db, err := sql.Open("mysql", ConnToString())
	if err != nil {
		panic(err.Error())
	}
	db.SetMaxOpenConns(20)
	DB = db
	// Does the database connection work?
	err = DB.Ping()
	if err != nil {
		fmt.Printf("Error could not ping database: %s\n", err.Error())
	}
	if DB == nil {
		log.Fatal("Database connection is nil")
	}
}

func Init() {
	fmt.Println("Connecting to database...")
	ConnectDB()
	// check if we can ping our DB
	err := DB.Ping()
	if err != nil {
		fmt.Printf("Error could not ping database: %s\n", err.Error())
	}
	if DB == nil {
		log.Fatal("Database connection is nil")
	}
	fmt.Println("Connected to database")
}

// package main

// import (
//     "log"
//     "net/http"
// )

// var db *sql.DB // Example database connection

// func main() {
//     // Initialize the database connection
//     var err error
//     db, err = sql.Open("mysql", "user:password@/dbname")
//     if err != nil {
//         log.Fatal(err)
//     }

//     // Your existing router setup
//     router := mux.NewRouter()
//     router.HandleFunc("/EventGET", events.EventGET).Methods("GET")
//     // Other routes

//     log.Fatal(http.ListenAndServe(":8080", router))
// }
