package users

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

func RegisterHandler(router *gin.Engine, db *sql.DB) {
	router.POST("/createTableUsers", HandlerCreateTableUsers(db))
	router.GET("/user", HandlerGetUserIDByName(db))
	router.POST("/user", HandlerInsertUser(db))
}
