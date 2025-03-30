package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/stainton/database/pkg/lottery/model"
)

// AddAnUser增加一个用户，成功返回userid
func (lc *LotteryClient) AddAnUser(name, telephone string) (userid int, err error) {
	// http客户端发起请求
	client := http.DefaultClient
	url := lc.getUrl(model.PATH_USER_ROOT)
	usr := &model.User{
		Name:      name,
		Telephone: telephone,
	}
	payload, err := json.Marshal(usr)
	if err != nil {
		lc.logger.Errorf("failed to marshal input: %v", err)
		return -1, err
	}
	response, err := client.Post(url, "application/json", bytes.NewReader(payload))
	if err != nil {
		lc.logger.Errorf("failed to reuqest: %v", err)
		return -1, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		lc.logger.Infof("request failed, status code: %d", response.StatusCode)
		return -1, err
	}
	buffer := make([]byte, 1024)
	n, err := response.Body.Read(buffer)
	if err != nil && err != io.EOF {
		lc.logger.Errorf("failed to read from body: %v, bytes: %d", err, n)
		return -1, err
	}
	err = json.Unmarshal(buffer[:n], usr)
	if err != nil {
		lc.logger.Errorf("failed to unmarshal: %v, bytes: %d", err, n)
		return -1, err
	}
	return usr.UserID, nil
}

// GetUsers查询用户信息，返回用户信息列表
func (lc *LotteryClient) GetUsers(usr *model.User, amount int) []*model.User {
	cli := http.DefaultClient
	url := lc.getUrl(model.PATH_USER_ROOT)
	if usr != nil {
		url = lc.withQuery(url, usr.Mapping())
		lc.logger.Infof("url %v", url)
	} else if amount > 0 {
		url = lc.withQuery(url, map[string]any{
			"amount": amount,
		})
		lc.logger.Infof("url %v", url)
	} else {
		lc.logger.Errorf("no valid amount: %v", amount)
		return nil
	}
	resp, err := cli.Get(url)
	if err != nil {
		lc.logger.Errorf("failed to get users: %v", err)
		return nil
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		lc.logger.Infof("request failed, status code: %d", resp.StatusCode)
		return nil
	}
	buffer, err := io.ReadAll(resp.Body)
	if err != nil {
		lc.logger.Errorf("failed to read from body: %v", err)
		return nil
	}

	var users []*model.User
	err = json.Unmarshal(buffer, &users)
	if err != nil {
		lc.logger.Errorf("failed to unmarshal users: %v", err)
		lc.logger.Debugf("response body is %v", bytes.NewBuffer(buffer).String())
		return nil
	}
	return users
}

// UpdateUser更新用户信息，以userid为准
func (lc *LotteryClient) UpdateUser(usr *model.User) error {
	if usr == nil || usr.UserID <= 0 {
		lc.logger.Errorf("valid userid is required for update: %v", usr)
		return errors.New("valid userid is required for update")
	}
	cli := http.DefaultClient
	url := lc.getUrl(model.PATH_USER_ROOT)
	buffer, err := json.Marshal(usr)
	if err != nil {
		lc.logger.Errorf("failed to marshal user: %v", err)
		return err
	}
	request, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(buffer))
	if err != nil {
		lc.logger.Errorf("failed to create request: %v", err)
		return err
	}
	response, err := cli.Do(request)
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

func (lc *LotteryClient) CreateUserTable() error {
	return lc.createUserTable(model.PATH_USER_TABLE_CREATE)
}
