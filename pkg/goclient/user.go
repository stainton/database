package goclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/stainton/database/pkg/model"
)

// GetUserIDByName 根据用户名称获取用户ID
func (c *Client) GetUserIDByName(name string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	url := fmt.Sprintf("%s/user?username=%s", c.baseURL(), name)
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		// create a new request failed
		return -1, errors.New("create request failed with " + err.Error())
	}
	response, err := c.client.Do(request)
	if err != nil {
		// request failed
		return -1, errors.New("request failed with " + err.Error())
	}
	defer response.Body.Close()
	user := model.User{}
	err = json.NewDecoder(response.Body).Decode(&user)
	if err != nil {
		// decode failed
		return -1, errors.New("decode response failed with " + err.Error())
	}
	if response.StatusCode != http.StatusOK {
		return -1, fmt.Errorf("get user failed: status code %d", response.StatusCode)
	}
	return user.Userid, nil
}

// RegisterUser 注册一个新用户
func (c *Client) RegisterUser(user *model.User) error {
	url := fmt.Sprintf("%s/user", c.baseURL())
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	userBytes, err := json.Marshal(user)
	if err != nil {
		return errors.New("parse payload failed with " + err.Error())
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(userBytes))
	if err != nil {
		return errors.New("create failed with " + err.Error())
	}
	response, err := c.client.Do(request)
	if err != nil {
		return errors.New("request failed with " + err.Error())
	}
	defer response.Body.Close()
	baseResponse := model.BaseReponsePayload{}
	err = json.NewDecoder(response.Body).Decode(&baseResponse)
	if response.StatusCode != http.StatusOK {
		return errors.New("create user failed with err:" + err.Error() + " and message" + baseResponse.Message)
	}
	return nil
}

// InitialUserTable 创建用户存储用户信息的表
func (c *Client) InitialUserTable() error {
	url := fmt.Sprintf("%s/createTableUsers", c.baseURL())
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return errors.New("create failed with " + err.Error())
	}
	response, err := c.client.Do(request)
	if err != nil {
		return errors.New("request failed with " + err.Error())
	}
	defer response.Body.Close()
	baseResponse := model.BaseReponsePayload{}
	err = json.NewDecoder(response.Body).Decode(&baseResponse)
	if response.StatusCode != http.StatusOK {
		return errors.New("create user failed with err:" + err.Error() + " and message" + baseResponse.Message)
	}
	return nil
}

func NewClient(remoteHost string, port int) *Client {
	return &Client{
		remoteHost: remoteHost,
		port:       port,
		client:     http.DefaultClient,
	}
}
