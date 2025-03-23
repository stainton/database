package order

import (
	"github.com/gin-gonic/gin"
	"github.com/stainton/database/internal/lottery/config"
	"github.com/stainton/logger"
)

func ChainMake(e *gin.Engine, l logger.Logger, rc *config.RuntimeConfig) {
	// 创建一个订单
	e.POST("/order", CreateOrderHandler(l, rc))
	// 数据库中创建orders表
	e.POST("/order/register", TableCreateHandler(l, rc))
	// 获取订单
	e.GET("/order", QueryOrderHandler(l, rc))
	// 更新order
	e.PUT("/order/:orderid", UpdateOrderHandler(l, rc))
}
