package cmd

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

type User struct {
	Userid   int64  `json:"userid"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

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
		username := ctx.Param("username")
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

func InsertUser(db *sql.DB, user *User) error {
	query := "INSERT INTO users (username, password, role) VALUES (?,?,?)"
	_, err := db.Exec(query, user.Username, user.Password, user.Role)
	return err
}

func HandlerInsertUser(db *sql.DB) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		user := User{}
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
