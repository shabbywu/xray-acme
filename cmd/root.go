// Package cmd 提供 Xray ACME 的命令行界面
// 使用 Cobra 实现命令行参数解析和子命令管理
package cmd

import (
	"context"
	"os"
	"xray-acme/cmd/api"

	"github.com/spf13/cobra"

	// 导入 Xray 所有内置功能
	_ "github.com/xtls/xray-core/main/distro/all"

	"xray-acme/config"
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

var (
	cfg *config.Config
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
	RunE:    run,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// 读取配置
		cfg = buildConfig()
		return nil
	},
}

// Execute 执行根命令
// 这是 CLI 的主入口点
func Execute() {
	if err := rootCmd.ExecuteContext(context.Background()); err != nil {
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

	rootCmd.AddCommand(api.Cmd)
}

// buildConfig 构建配置对象
// 命令行参数优先于环境变量
func buildConfig() *config.Config {
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
	return cfg
}

// getValueWithFallback 获取值，如果主值为空则使用备选值
func getValueWithFallback(primary, fallback string) string {
	if primary != "" {
		return primary
	}
	return fallback
}
