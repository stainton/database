package client

import (
	"context"
	"testing"

	"github.com/stainton/logger"
)

func TestAddAnUser(t *testing.T) {
	top := context.Background()
	l, _ := logger.NewLogger(top, &logger.Options{
		ServiceName:      "tester",
		Output:           "D:\\projects\\stainton\\database\\local_log\\test_log",
		MaxMessage:       1000,
		Threshold:        logger.MiB * 10,
		CompressInterval: 30,
	})
	cli := NewLotteryClient(top, l, "localhost", 8090)
	userid, err := cli.AddAnUser("harrgga", "165582452")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("add user success, userid:", userid)
}
