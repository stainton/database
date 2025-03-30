package client

import (
	"context"
	"fmt"
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
		// 无筛选无数量，返回nil
		{nil, -1, nil},
		// 空筛选无数量，返回nil
		{&model.User{}, -1, nil},
		// 有筛选无数量， 返回单个用户
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
		// 有筛选无数量， 返回单个用户
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
		// 有筛选有数量， 返回单个用户
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
		// 无筛选有数量， 返回单个用户
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

func TestUpdateUser(t *testing.T) {
	top := context.Background()
	l, _ := logger.NewLogger(top, &logger.Options{
		ServiceName:      "TestUpdateUser",
		Output:           "D:\\projects\\stainton\\database\\local_log\\test_log",
		MaxMessage:       1000,
		Threshold:        logger.MiB * 10,
		CompressInterval: 30,
	})
	cli := NewLotteryClient(top, l, "localhost", 8090)
	var exp []*model.User
	for i := range 20 {
		tc := &model.User{
			UserID:    i,
			Name:      fmt.Sprintf("user%d", i),
			Telephone: fmt.Sprintf("usr%d-tel", i),
		}
		us := cli.GetUsers(tc, -1)
		if len(us) == 0 {
			exp = nil
		} else {
			exp = []*model.User{tc}
		}
		err := cli.UpdateUser(tc)
		t.Logf("update result: %v", err)
		res := cli.GetUsers(tc, -1)
		if !reflect.DeepEqual(res, exp) {
			t.Errorf("Test case %d failed, expected %v, got %v", i, exp, res)
		}
	}
}

func TestCreateTable(t *testing.T) {
	top := context.Background()
	l, _ := logger.NewLogger(top, &logger.Options{
		ServiceName:      "TestUpdateUser",
		Output:           "D:\\projects\\stainton\\database\\local_log\\test_log",
		MaxMessage:       1000,
		Threshold:        logger.MiB * 10,
		CompressInterval: 30,
	})
	cli := NewLotteryClient(top, l, "localhost", 8090)
	if err := cli.CreateUserTable(); err != nil {
		t.Fatalf("create user table failed: %v", err)
	}
	if err := cli.CreateUserTable(); err == nil {
		t.Fatal("create a table repeatedly")
	}
}
