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
				Output:           "./logs",
				MaxMessage:       1000,
				Threshold:        logger.MiB * 10,
				CompressInterval: 30,
			})
			go func() {
				<-stop
				fmt.Println("Logger stoped...")
				wg.Done()
			}()
			sqlDB, err := db.NewDB(l, "root", "961110", "localhost", "lottery", 3306)
			if err != nil {
				cancel()
				wg.Wait()
				return
			}
			rc := &config.RuntimeConfig{DBhandler: sqlDB}
			router.POST("/user", user.RegisterUserHandler(l, rc))
			router.GET("/user", user.GetUserChain(l, rc)...)
			router.PUT("/user", user.UpdateUserHandler(l, rc))

			router.POST("/reward", reward.CreateRewardHandler(l, rc))

			router.POST("/order", order.CreateOrderHandler(l, rc))
			router.GET("/order", order.QueryOrderHandler(l, rc))
			router.PUT("/order/:orderid", order.QueryOrderHandler(l, rc))

			if err = router.Run(":8090"); err != nil {
				l.Fatalf("start server failed: %v", err)
			}
			cancel()
			wg.Wait()
			fmt.Println("Server stopped...")
		},
	}
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	return rootCmd
}
