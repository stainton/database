package order

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

// TODO: orderid不应该由用户创建，应该自增

// getQueryParams 解析URL中的查询变量，用于拼接查询语句，查询条件为Order的各个成员
func getQueryParams(c *gin.Context) string {
	elements := []string{}
	var v string
	var ok bool
	if v, ok = c.GetQuery("userid"); ok {
		elements = append(elements, fmt.Sprintf("userid = %s", v))
	}
	if v, ok = c.GetQuery("orderid"); ok {
		elements = append(elements, fmt.Sprintf("orderid = %s", v))
	}
	if v, ok = c.GetQuery("price"); ok {
		elements = append(elements, fmt.Sprintf("price = %s", v))
	}
	if v, ok = c.GetQuery("date"); ok {
		elements = append(elements, fmt.Sprintf("date = %s", v))
	}
	if v, ok = c.GetQuery("productid"); ok {
		elements = append(elements, fmt.Sprintf("productid = %s", v))
	}
	if v, ok = c.GetQuery("type"); ok {
		elements = append(elements, fmt.Sprintf("type = '%s'", v))
	}
	if len(elements) == 0 {
		return ""
	} else if len(elements) == 1 {
		return fmt.Sprintf("WHERE %s", elements[0])
	}
	return fmt.Sprintf("WHERE %s", strings.Join(elements, " AND "))
}

// QueryOrderHandler 获取满足条件的订单列表
func QueryOrderHandler(l logger.Logger, rc *config.RuntimeConfig) func(*gin.Context) {
	return func(c *gin.Context) {
		defer c.Abort()
		queryString := fmt.Sprintf("SELECT orderid,userid,productid,price,create_time,type FROM orders %s", getQueryParams(c))
		l.Info(queryString)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		rows, err := rc.DBhandler.QueryContext(ctx, queryString)
		if err != nil {
			l.Errorf("query order list failed: %v", err)
			c.JSON(http.StatusInternalServerError, common.ResponseTemplate{
				Code:    common.DB_QUERY_ERROR,
				Message: "query order specied failed.",
			})
			return
		}
		odrs := []model.Order{}
		for rows.Next() {
			o := model.Order{}
			err = rows.Scan(&o.OrderId, &o.UserId, &o.ProductId, &o.Price, &o.Date, &o.Type)
			if err != nil {
				l.Errorf("scan body from db failed: %v", err)
				continue
			}
			odrs = append(odrs, o)
		}
		c.JSON(http.StatusOK, odrs)
	}
}

// CreateOrderHandler 创建一个新的订单
func CreateOrderHandler(l logger.Logger, rc *config.RuntimeConfig) func(*gin.Context) {
	return func(c *gin.Context) {
		defer c.Abort()
		queryString := "INSERT INTO orders (userid,productid,price,create_time,type) VALUES (?,?,?,?,?)"
		odr := model.Order{}
		err := c.ShouldBindJSON(&odr)
		if err != nil {
			l.Errorf("unmarshal body failed: %v", err)
			c.JSON(http.StatusInternalServerError, common.ResponseTemplate{
				Code:    common.INTERNAL,
				Message: "create order failed.",
			})
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		res, err := rc.DBhandler.ExecContext(ctx, queryString, odr.UserId, odr.ProductId, odr.Price, odr.Date, odr.Type)
		if err != nil {
			l.Errorf("insert reward failed: %v", err)
			c.JSON(http.StatusInternalServerError, common.ResponseTemplate{
				Code:    common.DB_INSERT_ERROR,
				Message: "insert reward failed.",
			})
			return
		}
		id, err := res.LastInsertId()
		if err != nil {
			l.Errorf("get order failed: %v", err)
			id = -1
		}
		c.JSON(http.StatusOK, common.ResponseTemplate{
			Code:    int(id),
			Message: "order created successfully(code in body is the effected row).",
		})
	}
}

// UpdateOrderHandler 更新一个订单
func UpdateOrderHandler(l logger.Logger, rc *config.RuntimeConfig) func(*gin.Context) {
	return func(c *gin.Context) {
		orderid := c.Param("orderid")
		order := new(model.Order)
		err := c.ShouldBindJSON(order)
		if err != nil {
			c.JSON(http.StatusBadRequest, common.ResponseTemplate{
				Code:    common.DB_UPDATE_ERROR,
				Message: "update order failed.",
			})
			c.Abort()
			return
		}
		queryString := "UPDATE orders SET userid=?,productid=?,price=?,create_time=?,type=? WHERE orderid=?"
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		res, err := rc.DBhandler.ExecContext(ctx, queryString, order.UserId, order.ProductId, order.Price, order.Date, order.Type, orderid)
		if err != nil {
			l.Errorf("update order failed: %v", err)
			c.JSON(http.StatusInternalServerError, common.ResponseTemplate{
				Code:    common.DB_UPDATE_ERROR,
				Message: "update order failed.",
			})
			c.Abort()
			return
		}
		affect, err := res.RowsAffected()
		if err != nil || affect != 1 {
			l.Errorf("update order failed: %v, affected: %v", err, affect)
			c.JSON(http.StatusInternalServerError, common.ResponseTemplate{
				Code:    common.DB_UPDATE_ERROR,
				Message: "update order failed.",
			})
			c.Abort()
			return
		}
		c.JSON(http.StatusOK, common.ResponseTemplate{
			Code:    common.SUCCESS,
			Message: "order updated successfully.",
		})
	}
}

// TableCreateHandler 在数据库中创建orders表
func TableCreateHandler(l logger.Logger, rc *config.RuntimeConfig) func(*gin.Context) {
	return func(c *gin.Context) {
		defer c.Abort()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		db := rc.DBhandler
		queryString := `CREATE TABLE orders(
			orderid int NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT 'Primary Key',
			userid INT NOT NULL COMMENT 'User ID',
			productid INT NOT NULL COMMENT 'Product ID',
			price INT NOT NULL COMMENT 'Price',
			type ENUM('M', 'H') NOT NULL COMMENT 'model.Order Type',
			create_time DATETIME COMMENT 'Create Time',
			update_time DATETIME COMMENT 'Update Time'
		);`
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
