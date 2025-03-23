package reward

import (
	"github.com/gin-gonic/gin"
	"github.com/stainton/database/internal/lottery/config"
	"github.com/stainton/logger"
)

func ChainMake(e *gin.Engine, l logger.Logger, rc *config.RuntimeConfig) {
	e.GET("/reward", GetRewardHandler(l, rc))
	e.POST("/reward", CreateRewardHandler(l, rc))
	e.POST("/reward/register", TableCreateHandler(l, rc))
}
