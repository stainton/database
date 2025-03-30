package client

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/stainton/database/pkg/lottery/model"
)

func (lc *LotteryClient) CreateRewardTable() error {
	return lc.createUserTable(model.PATH_REWARD_TABLE_CREATE)
}

func (lc *LotteryClient) AddAReward(rwd *model.Reward) error {
	cli := http.DefaultClient
	url := lc.getUrl(model.PATH_REWARD_ROOT)
	buffer, err := json.Marshal(rwd)
	if err != nil {
		lc.logger.Errorf("marshal failed: %v", err)
		return err
	}
	response, err := cli.Post(url, "application/json", bytes.NewReader(buffer))
	if err != nil {
		lc.logger.Errorf("request failed: %v", err)
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		lc.logger.Infof("request failed, status code: %d", response.StatusCode)
		return err
	}
	return nil
}

func (lc *LotteryClient) GetRewards(qry *model.Reward) []*model.Reward {
	cli := http.DefaultClient
	url := lc.getUrl(model.PATH_REWARD_ROOT)
	url = lc.withQuery(url, qry.Mapping())
	response, err := cli.Get(url)
	if err != nil {
		lc.logger.Errorf("request failed: %v", err)
		return nil
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		lc.logger.Infof("request failed, status code: %d", response.StatusCode)
		return nil
	}
	var rewards []*model.Reward
	err = json.NewDecoder(response.Body).Decode(&rewards)
	if err != nil {
		lc.logger.Errorf("unmarshal failed: %v", err)
		return nil
	}
	return rewards
}
