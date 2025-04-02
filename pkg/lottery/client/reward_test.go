package client

import (
	"context"
	"fmt"
	"testing"

	"github.com/stainton/database/pkg/lottery/model"
	"github.com/stainton/logger"
)

func TestCreateRewardTable(t *testing.T) {
	top := context.Background()
	l, _ := logger.NewLogger(top, &logger.Options{
		ServiceName:      "TestUpdateUser",
		Output:           "D:\\projects\\stainton\\database\\local_log\\test_log",
		MaxMessage:       1000,
		Threshold:        logger.MiB * 10,
		CompressInterval: 30,
	})
	cli := NewLotteryClient(top, l, "localhost", 8090)
	err := cli.CreateRewardTable()
	if err != nil {
		t.Fatal(err)
	}
	err = cli.CreateRewardTable()
	if err == nil {
		t.Fatal(fmt.Errorf("CreateRewardTable success, expected error"))
	}
}

// 当前只支持date检索，且只返回号码
func TestAddAReward(t *testing.T) {
	top := context.Background()
	l, _ := logger.NewLogger(top, &logger.Options{
		ServiceName:      "TestUpdateUser",
		Output:           "D:\\projects\\stainton\\database\\local_log\\test_log",
		MaxMessage:       1000,
		Threshold:        logger.MiB * 10,
		CompressInterval: 30,
	})
	cli := NewLotteryClient(top, l, "localhost", 8090)
	err := cli.AddAReward(&model.Reward{
		ProductId: 39,
		Date:      "2025-04-02",
		Type:      "H",
	})
	if err != nil {
		t.Fatal(err)
	}
	if r := cli.GetRewards(&model.Reward{
		ProductId: 38,
		Date:      "2025-04-02",
		Type:      "H",
	}); r == -1 {
		// 预期找得到对应的号码，所以返回值不应该是-1
		t.Fatal(fmt.Errorf("no record found expected, rwd is %d", r))
	}
	if cli.GetRewards(&model.Reward{
		ProductId: 38,
		Date:      "2025-04-19",
		Type:      "H",
	}) != -1 {
		// 预期没有对应的号码，所以返回值应该是-1
		t.Fatal(fmt.Errorf("record found expected"))
	}
	if err = cli.AddAReward(&model.Reward{
		ProductId: 40,
		Date:      "2025-04-01",
		Type:      "H",
	}); err != nil {
		t.Fatal(err)
	}
	rws := cli.GetRewards(&model.Reward{
		ProductId: -1,
		Date:      "2025-04-01",
		Type:      "",
	})
	if rws != 40 {
		t.Fatal(fmt.Errorf("no record found"))
	}
}
