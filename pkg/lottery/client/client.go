package client

import (
	"context"
	"fmt"
	"strings"

	"github.com/stainton/logger"
)

type LotteryClient struct {
	hostname string
	port     int
	logger   logger.Logger
}

func (lc *LotteryClient) getUrl(pth string) string {
	return fmt.Sprintf("http://%s:%d/%s", lc.hostname, lc.port, pth)
}

func (lc *LotteryClient) withQuery(url string, mp map[string]any) string {
	if len(mp) == 0 {
		return url
	}
	qs := []string{}
	for k, v := range mp {
		qs = append(qs, fmt.Sprintf("%s=%v", k, v))
	}
	return fmt.Sprintf("%s?%s", url, strings.Join(qs, "&"))
}

func NewLotteryClient(ctx context.Context, l logger.Logger, hostname string, port int) *LotteryClient {
	return &LotteryClient{
		hostname: hostname,
		port:     port,
		logger:   l,
	}
}
