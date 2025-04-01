package client

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/stainton/database/pkg/lottery/model"
)

// e.POST("/order", CreateOrderHandler(l, rc))
// e.GET("/order", QueryOrderHandler(l, rc))
// e.PUT("/order/:orderid", UpdateOrderHandler(l, rc))

func (lc *LotteryClient) CreateOrderTable() error {
	return lc.createTable(model.PATH_ORDER_TABLE_CREATE)
}

func (lc *LotteryClient) UpdateOrder(ord *model.Order) error {
	cli := http.DefaultClient
	url := lc.withParams(lc.getUrl(model.PATH_ORDER_ROOT), ord.OrderId)
	buffer, err := json.Marshal(ord)
	if err != nil {
		lc.logger.Errorf("marshal order failed: %v", err)
		return err
	}
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(buffer))
	if err != nil {
		lc.logger.Errorf("failed to create request: %v", err)
		return err
	}
	res, err := cli.Do(req)
	if err != nil {
		lc.logger.Errorf("failed to update order: %v", err)
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		lc.logger.Errorf("update order failed, status code: %d", res.StatusCode)
		buffer, err := io.ReadAll(res.Body)
		lc.logger.Infof("response body: %v", bytes.NewBuffer(buffer).String())
		return err
	}
	return nil
}

func (lc *LotteryClient) AddAnOrder(ord *model.Order) error {
	// cli := http.DefaultClient
	// url := lc.getUrl(model.PATH_ORDER_ROOT)
	return nil
}
