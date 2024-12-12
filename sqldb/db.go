package sqldb

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

var (
	DB   *sql.DB
	once sync.Once
)

func InitDB() {
	once.Do(func() {
		var err error

		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
			os.Getenv("MYSQL_USER"),
			os.Getenv("MYSQL_PASSWORD"),
			os.Getenv("MYSQL_HOST"),
			os.Getenv("MYSQL_PORT"),
			os.Getenv("MYSQL_DATABASE"))
		DB, err = sql.Open("mysql", dsn)
		if err != nil {
			log.Fatalf("Error opening database: %v", err)
		}

		err = DB.Ping()
		if err != nil {
			log.Fatalf("Error connecting to database: %v", err)
		}

		log.Println("Database connection successfully established")
	})
}
