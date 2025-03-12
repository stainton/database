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
	"github.com/stainton/logger"
)

type Order struct {
	OrderId   int64  `json:"orderid"`
	UserId    int64  `json:"userid"`
	ProductId int64  `json:"productid"`
	Price     int64  `json:"price"`
	Date      string `json:"date"`
	Type      string `json:"type"`
}

func getQueryParams(c *gin.Context) string {
	elements := []string{}
	var v string
	var ok bool
	if v, ok = c.GetQuery("userid"); ok {
		elements = append(elements, fmt.Sprintf("userid = %s", v))
	}
	// 这里要根据tel/name来查到对应的userid
	// if v, ok = c.GetQuery("tel"); ok {
	// 	elements = append(elements, fmt.Sprintf("userid = %s", v))
	// }
	// if v, ok = c.GetQuery("name"); ok {
	// 	elements = append(elements, fmt.Sprintf("userid = %s", v))
	// }

	if v, ok = c.GetQuery("productid"); ok {
		elements = append(elements, fmt.Sprintf("productid = %s", v))
	}
	if v, ok = c.GetQuery("type"); ok {
		elements = append(elements, fmt.Sprintf("type = '%s'", v))
	}
	qp := strings.Join(elements, " AND ")
	return fmt.Sprintf("WHERE %s", qp)
}

func QueryOrderHandler(l logger.Logger, rc *config.RuntimeConfig) func(*gin.Context) {
	return func(c *gin.Context) {
		defer c.Abort()
		queryString := fmt.Sprintf("SELECT userid,productid,pay,create_time,type FROM orders %s", getQueryParams(c))
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
		odrs := []Order{}
		for rows.Next() {
			o := Order{}
			err = rows.Scan(&o.UserId, &o.ProductId, &o.Price, &o.Date, &o.Type)
			if err != nil {
				l.Errorf("scan body from db failed: %v", err)
				continue
			}
			odrs = append(odrs, o)
		}
		c.JSON(http.StatusOK, odrs)
	}
}

func CreateOrderHandler(l logger.Logger, rc *config.RuntimeConfig) func(*gin.Context) {
	return func(c *gin.Context) {
		defer c.Abort()
		queryString := "INSERT INTO orders (userid,productid,pay,create_time,type) VALUES (?,?,?,?,?)"
		odr := Order{}
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

func UpdateOrderHandler(l logger.Logger, rc *config.RuntimeConfig) func(*gin.Context) {
	return func(ctx *gin.Context) {}
}
