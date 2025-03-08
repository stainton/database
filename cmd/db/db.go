package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stainton/logger"
)

func NewDB(l logger.Logger, username, password, endpoint, database string, port int) (*sql.DB, error) {
	dbSourceString := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", username, password, endpoint, database)
	db, err := sql.Open("mysql", dbSourceString)
	if err != nil {
		l.Fatal("failed to open database")
		os.Exit(1)
	}
	return db, nil
}
