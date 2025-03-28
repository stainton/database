package user

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stainton/database/internal/lottery/common"
	"github.com/stainton/database/internal/lottery/config"
	"github.com/stainton/database/pkg/lottery/model"
	"github.com/stainton/logger"
)

// RegisterUserHandler 处理单个用户的注册
func RegisterUserHandler(l logger.Logger, rc *config.RuntimeConfig) func(*gin.Context) {
	return func(c *gin.Context) {
		db := rc.DBhandler
		usr := model.User{}
		if err := c.ShouldBindBodyWithJSON(&usr); err != nil {
			l.Errorf("unmarshal request body failed: %v", err)
			c.JSON(http.StatusBadRequest, common.ResponseTemplate{
				Code:    common.INVALID_PARAMS,
				Message: "invalid request body.",
			})
			return
		}
		queryString := "INSERT INTO users (name, telephone) VALUES (?,?)"
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		res, err := db.ExecContext(ctx, queryString, usr.Name, usr.Telephone)
		if err != nil {
			l.Errorf("insert user failed: %v", err)
			c.JSON(http.StatusInternalServerError, common.ResponseTemplate{
				Code:    common.DB_INSERT_ERROR,
				Message: "insert user failed.",
			})
			return
		}
		id, err := res.LastInsertId()
		if err != nil {
			l.Errorf("get userid failed: %v", err)
			id = -1
		}
		c.JSON(http.StatusOK, model.User{
			UserID:    int(id),
			Name:      usr.Name,
			Telephone: usr.Telephone,
		})
	}
}

// GetUserChain 处理用户的获取
func GetUserChain(l logger.Logger, rc *config.RuntimeConfig) gin.HandlersChain {
	return gin.HandlersChain{
		GetUserListHandler(l, rc),
		GetAnUserHandler(l, rc),
	}
}

// GetUserListHandler 获取用户列表
func GetUserListHandler(l logger.Logger, rc *config.RuntimeConfig) func(*gin.Context) {
	return func(c *gin.Context) {
		value, ok := c.GetQuery("batch")
		if !ok {
			c.Next()
			return
		}
		db := rc.DBhandler
		var queryString string
		if !strings.Contains(value, ",") {
			queryString = fmt.Sprintf("SELECT userid,name,telephone FROM users LIMIT %s", value)
		} else {
			vs := strings.Split(value, ",")
			queryString = fmt.Sprintf("SELECT userid,name,telephone FROM users LIMIT %s,%s", vs[0], vs[1])
		}
		l.Debugf("user list query is [%s]", queryString)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		rows, err := db.QueryContext(ctx, queryString)
		if err != nil {
			l.Errorf("query user list failed: %v", err)
			c.JSON(http.StatusInternalServerError, common.ResponseTemplate{
				Code:    common.DB_QUERY_ERROR,
				Message: "query user specied failed.",
			})
			return
		}

		usrs := []*model.User{}
		for rows.Next() {
			usr := model.User{}
			if err := rows.Scan(&usr.UserID, &usr.Name, &usr.Telephone); err != nil {
				l.Errorf("scan user info failed: %v", err)
				c.JSON(http.StatusInternalServerError, common.ResponseTemplate{
					Code:    common.DB_RESULT_SCAN,
					Message: "get user info failed.",
				})
				continue
			}
			usrs = append(usrs, &usr)
		}
		c.JSON(http.StatusOK, usrs)
		c.Abort()
	}
}

// GetAnUserHandler 获取单个用户
func GetAnUserHandler(l logger.Logger, rc *config.RuntimeConfig) func(*gin.Context) {
	return func(c *gin.Context) {
		db := rc.DBhandler
		var queryString, value string
		if value = c.Query("id"); value != "" {
			queryString = "SELECT userid,name,telephone FROM users WHERE userid = ?"
		} else if value = c.Query("name"); value != "" {
			queryString = "SELECT userid,name,telephone FROM users WHERE name = ?"
		} else if value = c.Query("tel"); value != "" {
			queryString = "SELECT userid,name,telephone FROM users WHERE telephone = ?"
		} else {
			l.Error("not params for querying.")
			c.JSON(http.StatusBadRequest, common.ResponseTemplate{
				Code:    common.INVALID_PARAMS,
				Message: "invalid request params.",
			})
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		row := db.QueryRowContext(ctx, queryString, value)
		usr := model.User{}
		if err := row.Scan(&usr.UserID, &usr.Name, &usr.Telephone); err != nil {
			l.Errorf("scan user info failed: %v", err)
			c.JSON(http.StatusInternalServerError, common.ResponseTemplate{
				Code:    common.DB_RESULT_SCAN,
				Message: "get user info failed.",
			})
			return
		}
		c.JSON(http.StatusOK, usr)
	}
}

// UpdateUserHandler 更新单个用户
func UpdateUserHandler(l logger.Logger, rc *config.RuntimeConfig) func(*gin.Context) {
	return func(c *gin.Context) {
		db := rc.DBhandler
		usr := model.User{}
		if err := c.ShouldBindBodyWithJSON(&usr); err != nil {
			l.Errorf("unmarshal request body failed: %v", err)
			c.JSON(http.StatusBadRequest, common.ResponseTemplate{
				Code:    common.INVALID_PARAMS,
				Message: "invalid request body.",
			})
			return
		}
		queryString := "UPDATE users SET name =?, telephone =? WHERE userid =?"
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		res, err := db.ExecContext(ctx, queryString, usr.Name, usr.Telephone, usr.UserID)
		if err != nil {
			l.Errorf("update user failed: %v", err)
			c.JSON(http.StatusInternalServerError, common.ResponseTemplate{
				Code:    common.DB_UPDATE_ERROR,
				Message: "update user failed.",
			})
			return
		}
		id, err := res.LastInsertId()
		if err != nil {
			l.Errorf("get userid failed: %v", err)
			id = -1
		}
		c.JSON(http.StatusOK, model.User{
			UserID:    int(id),
			Name:      usr.Name,
			Telephone: usr.Telephone,
		})
	}
}

func TableCreateHandler(l logger.Logger, rc *config.RuntimeConfig) func(*gin.Context) {
	return func(c *gin.Context) {
		defer c.Abort()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		db := rc.DBhandler
		queryString := `CREATE TABLE users(
			userid int NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT 'Primary Key',
			create_time DATETIME COMMENT 'Create Time',
			update_time DATETIME COMMENT 'Update Time',
			telephone VARCHAR(128) UNIQUE NOT NULL COMMENT 'Telephone Number',
			name VARCHAR(128));`
		_, err := db.ExecContext(ctx, queryString)
		if err != nil {
			l.Errorf("create table failed: %v", err)
			c.JSON(http.StatusInternalServerError, common.ResponseTemplate{
				Code:    common.DB_CONNECT,
				Message: "create table failed.",
			})
			return
		}
		c.JSON(http.StatusOK, common.ResponseTemplate{
			Code:    common.SUCCESS,
			Message: "create table successfully.",
		})
	}
}
