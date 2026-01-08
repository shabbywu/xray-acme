// Package cmd 提供 Xray ACME 的命令行界面
// 使用 Cobra 实现命令行参数解析和子命令管理
package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"

	// 导入 Xray 所有内置功能
	_ "github.com/xtls/xray-core/main/distro/all"

	"xray-acme/acme"
	"xray-acme/config"
	"xray-acme/launcher"
)

// Version 程序版本号
var Version = "1.0.0"

// 命令行参数
var (
	cfgFile string // Xray 配置文件路径
	dpID    string // DNSPod API ID
	dpToken string // DNSPod API Token
	email   string // Let's Encrypt 注册邮箱
)

// rootCmd 根命令
var rootCmd = &cobra.Command{
	Use:   "xray-acme",
	Short: "自动管理 TLS 证书的 Xray 启动器",
	Long: `Xray ACME - 自动管理 TLS 证书的 Xray 启动器

使用 CertMagic 和 DNSPod DNS-01 挑战自动申请和续期 Let's Encrypt 证书，
并将证书自动注入到 Xray 配置中，实现全自动的 TLS 证书管理。

示例:
  # 使用命令行参数
  xray-acme --config config.json --dp-id YOUR_ID --dp-token YOUR_TOKEN --email your@email.com

  # 使用环境变量
  export DP_ID="your-id"
  export DP_TOKEN="your-token"
  export EMAIL="your@email.com"
  xray-acme`,
	Version: Version,
	RunE:    runRoot,
}

// Execute 执行根命令
// 这是 CLI 的主入口点
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// 绑定命令行参数
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "config.json",
		"Xray 配置文件路径")
	rootCmd.PersistentFlags().StringVar(&dpID, "dp-id", "",
		"DNSPod API ID (也可通过环境变量 DP_ID 设置)")
	rootCmd.PersistentFlags().StringVar(&dpToken, "dp-token", "",
		"DNSPod API Token (也可通过环境变量 DP_TOKEN 设置)")
	rootCmd.PersistentFlags().StringVar(&email, "email", "",
		"Let's Encrypt 注册邮箱 (也可通过环境变量 EMAIL 设置)")
}

// runRoot 根命令的执行逻辑
func runRoot(cmd *cobra.Command, args []string) error {
	log.Printf("Xray ACME v%s 启动中...", Version)

	// 1. 构建配置（命令行参数优先于环境变量）
	cfg, err := buildConfig()
	if err != nil {
		return fmt.Errorf("配置加载失败: %w", err)
	}

	// 2. 创建证书管理器
	certManager := acme.NewManager(cfg)

	// 3. 创建并运行启动器
	xrayLauncher := launcher.New(cfg, certManager)

	ctx := context.Background()
	if err := xrayLauncher.Run(ctx); err != nil {
		return fmt.Errorf("启动失败: %w", err)
	}

	log.Println("Xray ACME 已退出")
	return nil
}

// buildConfig 构建配置对象
// 命令行参数优先于环境变量
func buildConfig() (*config.Config, error) {
	cfg := &config.Config{
		DNS: config.DNSConfig{
			ProviderID:    getValueWithFallback(dpID, os.Getenv(config.EnvDNSPodID)),
			ProviderToken: getValueWithFallback(dpToken, os.Getenv(config.EnvDNSPodToken)),
		},
		Cert: config.CertConfig{
			Email: getValueWithFallback(email, os.Getenv(config.EnvEmail)),
		},
		Xray: config.XrayConfig{
			ConfigPath: cfgFile,
		},
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// getValueWithFallback 获取值，如果主值为空则使用备选值
func getValueWithFallback(primary, fallback string) string {
	if primary != "" {
		return primary
	}
	return fallback
}
