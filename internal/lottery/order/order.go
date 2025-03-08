package order

import (
	"github.com/gin-gonic/gin"
	"github.com/stainton/database/internal/lottery/config"
	"github.com/stainton/logger"
)

type Order struct {
	OrderId   int64  `json:"orderid"`
	UserId    int64  `json:"userid"`
	ProductId int64  `json:"productid"`
	Price     int64  `json:"price"`
	Date      string `json:"date"`
	Type      int    `json:"type"`
}

func QueryOrderHandler(l logger.Logger, rc *config.RuntimeConfig) func(*gin.Context) {
	return func(ctx *gin.Context) {}
}

func CreateOrderHandler(l logger.Logger, rc *config.RuntimeConfig) func(*gin.Context) {
	return func(ctx *gin.Context) {}
}

func UpdateOrderHandler(l logger.Logger, rc *config.RuntimeConfig) func(*gin.Context) {
	return func(ctx *gin.Context) {}
}
