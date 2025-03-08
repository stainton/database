package reward

import (
	"github.com/gin-gonic/gin"
	"github.com/stainton/database/internal/lottery/config"
	"github.com/stainton/logger"
)

type Reward struct {
	ProductId int64  `json:"productid"`
	Date      string `json:"date"`
}

func GetRewardHandler(l logger.Logger, rc *config.RuntimeConfig) func(*gin.Context) {
	return func(ctx *gin.Context) {
	}
}

func CreateRewardHandler(l logger.Logger, rc *config.RuntimeConfig) func(*gin.Context) {
	return func(ctx *gin.Context) {
	}
}
