package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func InitConnection() *sql.DB {

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, password, host, port, dbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connecting to MySQL: %v", err)
	}

	log.Println("¡Connecting MySQL success! 🚀")
	return db
}
