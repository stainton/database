package cmd

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"

	_ "github.com/go-sql-driver/mysql"
)

func NewDbServer() *cobra.Command {
	var password string
	var dbName string
	var service string
	dbCommand := &cobra.Command{
		Use:   "database",
		Short: "A brief description of your application",
		Long: `A longer description that spans multiple lines and likely contains
	examples and usage of using your application. For example:

	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a Cobra application.`,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		Run: func(cmd *cobra.Command, args []string) {
			db, err := GetDatabase(service, dbName, password)
			if err != nil {
				os.Exit(1)
			}
			RunServer(db)
		},
	}
	dbCommand.Flags().StringVarP(&password, "password", "p", "961110", "Password to Mysql database")
	dbCommand.Flags().StringVarP(&dbName, "database", "d", "testdb", "Name of Mysql database")
	dbCommand.Flags().StringVarP(&service, "service", "s", "0.0.0.0", "database address")
	return dbCommand
}

func GetDatabase(service, database, password string) (*sql.DB, error) {
	connectString := fmt.Sprintf("root:%s@tcp(%s:3306)/%s", password, service, database)
	db, err := sql.Open("mysql", connectString)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func RunServer(db *sql.DB) error {
	router := gin.Default()
	router.POST("/createTableUsers", HandlerCreateTableUsers(db))
	router.GET("/user/:username", HandlerGetUserIDByName(db))
	router.POST("/user", HandlerInsertUser(db))
	return router.Run(":8080")
}
