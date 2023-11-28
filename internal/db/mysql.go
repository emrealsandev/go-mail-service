package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"log"
)

var DB *sql.DB

func InitDB() {
	dbConfig := viper.GetStringMapString("database")

	dataSourceName := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		dbConfig["username"],
		dbConfig["password"],
		dbConfig["host"],
		dbConfig["port"],
		dbConfig["dbname"],
	)

	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	DB = db
}
