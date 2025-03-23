package user

import (
	"github.com/gin-gonic/gin"
	"github.com/stainton/database/internal/lottery/config"
	"github.com/stainton/logger"
)

func ChainMake(e *gin.Engine, l logger.Logger, rc *config.RuntimeConfig) {
	// 新增一个用户
	e.POST("/user", RegisterUserHandler(l, rc))
	// 数据库中创建user表
	e.POST("/user/register", TableCreateHandler(l, rc))
	// 获取用户信息
	e.GET("/user", GetUserChain(l, rc)...)
	// 更新用户信息
	e.PUT("/user", UpdateUserHandler(l, rc))
}
