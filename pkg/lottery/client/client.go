package client

import (
	"context"
	"fmt"

	"github.com/stainton/logger"
)

type LotteryStruct struct {
	hostname string
	port     int
	logger   logger.Logger
}

func (s *LotteryStruct) getUrl(pth string) string {
	return fmt.Sprintf("http://%s:%d/%s", s.hostname, s.port, pth)
}

func NewLotteryClient(ctx context.Context, l logger.Logger, hostname string, port int) *LotteryStruct {
	return &LotteryStruct{
		hostname: hostname,
		port:     port,
		logger:   l,
	}
}
