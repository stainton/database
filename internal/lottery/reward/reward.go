package reward

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stainton/database/internal/lottery/common"
	"github.com/stainton/database/internal/lottery/config"
	"github.com/stainton/database/pkg/lottery/model"
	"github.com/stainton/logger"
	"golang.org/x/net/context"
)

// 应该返回实际的购买情况
func GetRewardHandler(l logger.Logger, rc *config.RuntimeConfig) func(*gin.Context) {
	return func(c *gin.Context) {
		d, exist := c.GetQuery("date")
		if !exist {
			l.Infof("parameter 'date' not found")
			c.Next()
			return
		}
		defer c.Abort()
		l.Infof("query reward of date %s", d)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		queryString := "SELECT productid,date FROM rewards WHERE date = ?"
		row := rc.DBhandler.QueryRowContext(ctx, queryString, d)
		// l.Infof("reward is %s,%s", )
		productid := -1
		if err := row.Scan(&productid, &d); err != nil || productid == -1 {
			l.Errorf("scan reward from result in %s failed: %v, productid: %d", d, err, productid)
			c.JSON(http.StatusInternalServerError, common.ResponseTemplate{
				Code:    common.DB_RESULT_SCAN,
				Message: "get reward failed.",
			})
			return
		}
		c.JSON(http.StatusOK, common.ResponseTemplate{
			Code:    productid,
			Message: "code in body is the reward.",
		})
	}
}

func CreateRewardHandler(l logger.Logger, rc *config.RuntimeConfig) func(*gin.Context) {
	return func(c *gin.Context) {
		defer c.Abort()
		reward := model.Reward{}
		err := c.ShouldBindBodyWithJSON(&reward)
		if err != nil {
			l.Errorf("unmarshal reward from body failed : %v", err)
			c.JSON(http.StatusBadRequest, common.ResponseTemplate{
				Code:    common.INVALID_PARAMS,
				Message: "invalid request body.",
			})
			return
		}
		_, err = time.Parse("2006-01-02", reward.Date)
		if err != nil {
			l.Errorf("unmarshal reward from body failed : %v", err)
			c.JSON(http.StatusBadRequest, common.ResponseTemplate{
				Code:    common.INVALID_PARAMS,
				Message: "invalid request body.",
			})
			return
		}
		queryString := "INSERT INTO rewards (productid, date) VALUES (?, ?)"
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		res, err := rc.DBhandler.ExecContext(ctx, queryString, reward.ProductId, reward.Date)
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
			l.Errorf("get userid failed: %v", err)
			id = -1
		}
		c.JSON(http.StatusOK, common.ResponseTemplate{
			Code:    int(id),
			Message: "reward created successfully(code in body is the effected row).",
		})
	}
}

func TableCreateHandler(l logger.Logger, rc *config.RuntimeConfig) func(*gin.Context) {
	return func(c *gin.Context) {
		defer c.Abort()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		db := rc.DBhandler
		queryString := `CREATE TABLE rewards(
			id int NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT 'Primary Key',
			productid INT NOT NULL COMMENT 'Product ID',
			type ENUM('M', 'H') NOT NULL COMMENT 'Rward type',
			date DATETIME COMMENT 'Date'
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
