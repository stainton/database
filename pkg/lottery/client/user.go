package client

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/stainton/database/pkg/lottery/model"
)

func (c *LotteryStruct) AddAnUser(name, telephone string) (userid int, err error) {
	// http客户端发起请求
	client := http.DefaultClient
	url := c.getUrl(model.PATH_USER_ROOT)
	usr := &model.User{
		Name:      name,
		Telephone: telephone,
	}
	payload, err := json.Marshal(usr)
	if err != nil {
		c.logger.Errorf("failed to marshal input: %v", err)
		return -1, err
	}
	response, err := client.Post(url, "application/json", bytes.NewReader(payload))
	if err != nil {
		c.logger.Errorf("failed to reuqest: %v", err)
		return -1, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		c.logger.Infof("request failed, status code: %d", response.StatusCode)
		return -1, err
	}
	buffer := make([]byte, 1024)
	n, err := response.Body.Read(buffer)
	if err != nil && err != io.EOF {
		c.logger.Errorf("failed to read from body: %v, bytes: %d", err, n)
		return -1, err
	}
	err = json.Unmarshal(buffer[:n], usr)
	if err != nil {
		c.logger.Errorf("failed to unmarshal: %v, bytes: %d", err, n)
		return -1, err
	}
	return usr.UserID, nil
}
