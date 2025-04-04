/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/stainton/database/cmd/db"
	"github.com/stainton/database/internal/lottery/config"
	"github.com/stainton/database/internal/lottery/order"
	"github.com/stainton/database/internal/lottery/reward"
	"github.com/stainton/database/internal/lottery/user"
	"github.com/stainton/logger"
)

func NewDatabaseAPICmd() *cobra.Command {
	// new和直接创建结构体有什么不同
	cfg := new(db.DBConfig)
	// rootCmd represents the base command when called without any subcommands
	rootCmd := &cobra.Command{
		Use:   "database",
		Short: "A brief description of your application",
		Long: `A longer description that spans multiple lines and likely contains
	examples and usage of using your application. For example:

	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a Cobra application.`,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		Run: func(cmd *cobra.Command, args []string) {
			router := gin.Default()
			top, cancel := context.WithCancel(context.Background())
			wg := sync.WaitGroup{}

			wg.Add(1)
			l, stop := logger.NewLogger(top, &logger.Options{
				ServiceName:      "database",
				Output:           cfg.GetLogDictionary(),
				MaxMessage:       1000,
				Threshold:        logger.MiB * 10,
				CompressInterval: 30,
			})
			go func() {
				<-stop
				fmt.Println("Logger stoped...")
				wg.Done()
			}()
			// sqlDB, err := db.NewDB(l, "root", "961110", "localhost", "lottery", 13306)
			sqlDB, err := cfg.NewDB()
			// TODO: 需要检测数据库是否可以联通
			if err != nil {
				l.Errorf("Error creating database: %v", err)
				cancel()
				wg.Wait()
				return
			}
			rc := &config.RuntimeConfig{DBhandler: sqlDB}
			user.ChainMake(router, l, rc)
			reward.ChainMake(router, l, rc)
			order.ChainMake(router, l, rc)

			listenAddress := cfg.GetListeningAddresses()
			l.Infof("listening on %v", listenAddress)
			if err = router.Run(listenAddress); err != nil {
				l.Fatalf("start server failed: %v", err)
			}
			cancel()
			wg.Wait()
			fmt.Println("Server stopped...")
		},
	}
	cfg.AddFlags(rootCmd)
	return rootCmd
}
