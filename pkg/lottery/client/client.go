package client

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/stainton/logger"
)

type LotteryClient struct {
	hostname string
	port     int
	logger   logger.Logger
}

func (lc *LotteryClient) join(sep string, qrys []any) string {
	res := []string{}
	for _, v := range qrys {
		res = append(res, fmt.Sprintf("%v", v))
	}
	return strings.Join(res, sep)
}

func (lc *LotteryClient) getUrl(pth string) string {
	return fmt.Sprintf("http://%s:%d/%s", lc.hostname, lc.port, pth)
}

func (lc *LotteryClient) withParams(url string, params ...any) string {
	if len(params) == 0 {
		return url
	}
	return fmt.Sprintf("%s/%s", url, lc.join("/", params))
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

func (lc *LotteryClient) createTable(pth string) error {
	cli := http.DefaultClient
	url := lc.getUrl(pth)
	response, err := cli.Post(url, "application/json", nil)
	if err != nil {
		lc.logger.Errorf("request failed: %v", err)
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		lc.logger.Infof("response status code: %v", response.StatusCode)
		return fmt.Errorf("response status code: %v", response.StatusCode)
	}
	return nil
}

func NewLotteryClient(ctx context.Context, l logger.Logger, hostname string, port int) *LotteryClient {
	return &LotteryClient{
		hostname: hostname,
		port:     port,
		logger:   l,
	}
}
