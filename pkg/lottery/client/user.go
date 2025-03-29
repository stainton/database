package client

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/stainton/database/pkg/lottery/model"
)

func (c *LotteryClient) AddAnUser(name, telephone string) (userid int, err error) {
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

func (c *LotteryClient) GetUsers(usr *model.User, amount int) []*model.User {
	cli := http.DefaultClient
	url := c.getUrl(model.PATH_USER_ROOT)
	if usr != nil {
		url = c.withQuery(url, usr.Mapping())
		c.logger.Infof("url %v", url)
	} else if amount > 0 {
		url = c.withQuery(url, map[string]any{
			"amount": amount,
		})
		c.logger.Infof("url %v", url)
	} else {
		c.logger.Errorf("no valid amount: %v", amount)
		return nil
	}
	resp, err := cli.Get(url)
	if err != nil {
		c.logger.Errorf("failed to get users: %v", err)
		return nil
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		c.logger.Infof("request failed, status code: %d", resp.StatusCode)
		return nil
	}
	buffer, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Errorf("failed to read from body: %v", err)
		return nil
	}

	var users []*model.User
	err = json.Unmarshal(buffer, &users)
	if err != nil {
		c.logger.Errorf("failed to unmarshal users: %v", err)
		c.logger.Debugf("response body is %v", bytes.NewBuffer(buffer).String())
		return nil
	}
	return users
}
