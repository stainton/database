package users

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

// GetUserIDByName 通过用户名称获取用户ID
func GetUserIDByName(db *sql.DB, username string) (int64, error) {
	query := "SELECT userid FROM users WHERE username =?"
	row := db.QueryRow(query, username)

	var userID int64
	err := row.Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func HandlerGetUserIDByName(db *sql.DB) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		username := ctx.Query("username")
		if username == "" {
			ctx.JSON(400, gin.H{"error": "username is required"})
			return
		}
		userID, err := GetUserIDByName(db, username)
		if err == sql.ErrNoRows {
			ctx.JSON(404, gin.H{"error": "User not found"})
			return
		}
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(200, gin.H{"userid": userID})
	}
}
