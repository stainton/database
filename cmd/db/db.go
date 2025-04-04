package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"
	"github.com/stainton/logger"
)

const (
	ENV_LOG_PATH = "ENV_LOG_PATH"
)

type DBConfig struct {
	Server        string
	Port          string
	User          string
	Password      string
	DatebaseName  string
	ListeningAddr string
	ListeningPort string
	LogDictionary string
}

func (cfg *DBConfig) AddFlags(cCmd *cobra.Command) {
	cCmd.Flags().StringVar(&cfg.Server, "db-server", "localhost", "addrress of the database")
	cCmd.Flags().StringVar(&cfg.Port, "db-port", "3306", "port of the database")
	cCmd.Flags().StringVar(&cfg.User, "db-user", "root", "user of the database")
	cCmd.Flags().StringVar(&cfg.Password, "db-passwd", "961110", "password of the database")
	cCmd.Flags().StringVar(&cfg.DatebaseName, "db-name", "test-db", "name of the database")
	cCmd.Flags().StringVar(&cfg.ListeningAddr, "addr", "localhost", "address the database api is listening on")
	cCmd.Flags().StringVar(&cfg.ListeningPort, "port", "8090", "port the database api is listening on")
	cCmd.Flags().StringVar(&cfg.LogDictionary, "log-output", "", "log output")
}

func (cfg *DBConfig) getEnvWithDefault(value, env, defValue string) string {
	if value != "" {
		return value
	}
	if v := os.Getenv(env); v != "" {
		return v
	}
	return defValue
}

func (cfg *DBConfig) GetLogDictionary() string {
	return cfg.getEnvWithDefault(cfg.LogDictionary, ENV_LOG_PATH, "./logs")
}

func (cfg *DBConfig) GetListeningAddresses() string {
	return fmt.Sprintf("%s:%s", cfg.ListeningAddr, cfg.ListeningPort)
}

func (cfg *DBConfig) NewDB() (*sql.DB, error) {
	dbSourceString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", cfg.User, cfg.Password, cfg.Server, cfg.Port, cfg.DatebaseName)
	return sql.Open("mysql", dbSourceString)
}

func NewDB(l logger.Logger, username, password, endpoint, database string, port int) (*sql.DB, error) {
	dbSourceString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", username, password, endpoint, port, database)
	db, err := sql.Open("mysql", dbSourceString)
	if err != nil {
		l.Fatal("failed to open database")
		os.Exit(1)
	}
	return db, nil
}
