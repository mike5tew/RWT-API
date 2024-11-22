package sqldb

import (
	"database/sql"
	"log"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

var (
	DB   *sql.DB
	once sync.Once
)

func InitDB(dataSourceName string) {
	once.Do(func() {
		var err error
		DB, err = sql.Open("mysql", dataSourceName)
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
