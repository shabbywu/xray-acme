// Xray ACME - 自动管理 TLS 证书的 Xray 启动器
//
// 本程序使用 CertMagic 和 DNSPod DNS-01 挑战自动申请和续期 Let's Encrypt 证书，
// 并将证书自动注入到 Xray 配置中，实现全自动的 TLS 证书管理。
//
// 使用方式:
//
//	# 使用命令行参数
//	xray-acme --config config.json --dp-id YOUR_ID --dp-token YOUR_TOKEN --email your@email.com
//
//	# 使用环境变量
//	export DP_ID="your-id"
//	export DP_TOKEN="your-token"
//	export EMAIL="your@email.com"
//	xray-acme
//
// 查看帮助:
//
//	xray-acme --help
package main

import "xray-acme/cmd"

func main() {
	cmd.Execute()
}
