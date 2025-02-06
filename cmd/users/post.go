package users

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/stainton/database/pkg/model"
)

// CreateTableUsers 创建用于存储用户信息的数据表
func CreateTableUsers(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS users (
        userid INT AUTO_INCREMENT PRIMARY KEY,
        username VARCHAR(255) UNIQUE NOT NULL,
        password VARCHAR(255) NOT NULL,
        role VARCHAR(255) NOT NULL
    );`
	_, err := db.Exec(query)
	return err
}

func HandlerCreateTableUsers(db *sql.DB) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		err := CreateTableUsers(db)
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(200, gin.H{"message": "Table users created successfully"})
	}
}

// InsertUser 往users数据表中插入一个新用户
func InsertUser(db *sql.DB, user *model.User) error {
	query := "INSERT INTO users (username, password, role) VALUES (?,?,?)"
	_, err := db.Exec(query, user.Username, user.Password, user.Role)
	return err
}

func HandlerInsertUser(db *sql.DB) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		user := model.User{}
		err := ctx.ShouldBindJSON(&user)
		if err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}
		err = InsertUser(db, &user)
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(200, gin.H{"message": "User inserted successfully"})
	}
}
