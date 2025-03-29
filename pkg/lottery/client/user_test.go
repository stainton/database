package client

import (
	"context"
	"reflect"
	"testing"

	"github.com/stainton/database/pkg/lottery/model"
	"github.com/stainton/logger"
)

func TestAddAnUser(t *testing.T) {
	top := context.Background()
	l, _ := logger.NewLogger(top, &logger.Options{
		ServiceName:      "TestAddAnUser",
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

func TestGetUser(t *testing.T) {
	top := context.Background()
	l, _ := logger.NewLogger(top, &logger.Options{
		ServiceName:      "TestGetUser",
		Output:           "D:\\projects\\stainton\\database\\local_log\\test_log",
		MaxMessage:       1000,
		Threshold:        logger.MiB * 10,
		CompressInterval: 30,
	})
	cli := NewLotteryClient(top, l, "localhost", 8090)
	tcs := [][]any{
		{nil, -1, nil},
		{&model.User{}, -1, nil},
		{&model.User{
			UserID:    1,
			Name:      "harrgga",
			Telephone: "165582452",
		}, -1, []*model.User{
			{
				UserID:    1,
				Name:      "masha",
				Telephone: "18149630903",
			},
		}},
		{&model.User{
			UserID:    -1,
			Name:      "masha",
			Telephone: "165582452",
		}, -1, []*model.User{
			{
				UserID:    1,
				Name:      "masha",
				Telephone: "18149630903",
			},
		}},
		{&model.User{
			UserID:    -1,
			Name:      "",
			Telephone: "18149630903",
		}, 5, []*model.User{
			{
				UserID:    1,
				Name:      "masha",
				Telephone: "18149630903",
			},
		}},
		{nil, 3, []*model.User{
			{
				UserID:    1,
				Name:      "masha",
				Telephone: "18149630903",
			},
			{
				UserID:    2,
				Name:      "air man",
				Telephone: "7506200",
			},
			{
				UserID:    6,
				Name:      "airman",
				Telephone: "7506201",
			},
		}},
	}
	for i, tc := range tcs {
		var p1 *model.User = nil
		if tc[0] != nil {
			p1 = tc[0].(*model.User)
		}
		var exp []*model.User = nil
		if tc[2] != nil {
			exp = tc[2].([]*model.User)
		}
		user := cli.GetUsers(p1, tc[1].(int))
		if !reflect.DeepEqual(user, exp) {
			t.Errorf("Test case %d failed, expected %v, got %v", i, exp, user)
		}
	}
}
