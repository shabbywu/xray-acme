package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"xray-acme/acme"
	"xray-acme/launcher"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run xray-core",
	Short: "Run Xray with config, the default command",
	Long:  "Run Xray with config, the default command.",
	RunE:  run,
}

// run xray-core 执行逻辑
func run(cmd *cobra.Command, args []string) error {
	log.Printf("Xray ACME v%s 启动中...", Version)

	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("配置加载失败: %w", err)
	}
	
	// 1. 创建证书管理器
	certManager := acme.NewManager(cfg)

	// 2. 创建并运行启动器
	xrayLauncher := launcher.New(cfg, certManager)

	if err := xrayLauncher.Run(cmd.Context()); err != nil {
		return fmt.Errorf("启动失败: %w", err)
	}

	log.Println("Xray ACME 已退出")
	return nil
}

func init() {
	rootCmd.AddCommand(runCmd)
}
