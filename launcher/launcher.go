// Package launcher 提供 Xray 服务的启动和管理功能
// 负责配置解析、证书注入和服务生命周期管理
package launcher

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"runtime/debug"
	"strings"
	"syscall"

	"github.com/xtls/xray-core/common/errors"
	"github.com/xtls/xray-core/core"
	"github.com/xtls/xray-core/infra/conf"
	"github.com/xtls/xray-core/infra/conf/serial"
	"github.com/xtls/xray-core/main/confloader"

	"xray-acme/acme"
	"xray-acme/config"
)

// Launcher Xray 服务启动器
// 负责加载配置、管理证书并启动 Xray 服务
type Launcher struct {
	// config 应用程序配置
	config *config.Config

	// certManager 证书管理器
	certManager *acme.Manager

	// xrayConfig Xray 配置对象
	xrayConfig *conf.Config

	// server Xray 服务实例
	server *core.Instance
}

// New 创建新的 Launcher 实例
//
// 参数:
//   - cfg: 应用程序配置
//   - certMgr: 证书管理器
//
// 返回:
//   - *Launcher: Launcher 实例
func New(cfg *config.Config, certMgr *acme.Manager) *Launcher {
	return &Launcher{
		config:      cfg,
		certManager: certMgr,
	}
}

// LoadConfig 加载并解析 Xray 配置文件
//
// 返回:
//   - error: 如果加载或解析失败则返回错误
func (l *Launcher) LoadConfig() error {
	configPath := l.config.Xray.ConfigPath
	log.Printf("正在加载 Xray 配置: %s", configPath)

	// 读取配置文件
	reader, err := confloader.LoadConfig(configPath)
	if err != nil {
		return errors.New("读取配置文件失败: ", configPath, " - ", err)
	}

	// 解析 JSON 配置
	xrayConfig, err := serial.DecodeJSONConfig(reader)
	if err != nil {
		return errors.New("解析配置文件失败: ", configPath, " - ", err)
	}

	l.xrayConfig = xrayConfig
	log.Println("Xray 配置加载成功")
	return nil
}

// ExtractDomains 从 Xray 配置中提取需要证书的域名列表
//
// 返回:
//   - []string: 域名列表
func (l *Launcher) ExtractDomains() []string {
	var domains []string

	if l.xrayConfig == nil {
		return domains
	}

	for _, inbound := range l.xrayConfig.InboundConfigs {
		if inbound.StreamSetting == nil || inbound.StreamSetting.TLSSettings == nil {
			continue
		}
		tlsSettings := inbound.StreamSetting.TLSSettings
		if tlsSettings.ServerName != "" {
			domains = append(domains, tlsSettings.ServerName)
		}
	}

	return domains
}

// InjectCertificates 将证书注入到 Xray 配置中
//
// 返回:
//   - error: 如果获取或注入证书失败则返回错误
func (l *Launcher) InjectCertificates() error {
	if l.xrayConfig == nil {
		return fmt.Errorf("Xray 配置未加载")
	}

	log.Println("正在注入 TLS 证书到配置...")

	for _, inbound := range l.xrayConfig.InboundConfigs {
		if inbound.StreamSetting == nil || inbound.StreamSetting.TLSSettings == nil {
			continue
		}

		tlsSettings := inbound.StreamSetting.TLSSettings
		if tlsSettings.ServerName == "" {
			continue
		}

		domain := tlsSettings.ServerName
		log.Printf("正在为域名 %s 注入证书", domain)

		// 获取证书密钥对
		certKeyPair, err := l.certManager.GetCertKeyPair(domain)
		if err != nil {
			return fmt.Errorf("获取域名 %s 的证书失败: %w", domain, err)
		}

		// 注入到配置中
		tlsSettings.Certs = []*conf.TLSCertConfig{
			{
				CertStr: strings.Split(certKeyPair.CertPEM, "\n"),
				KeyStr:  strings.Split(certKeyPair.KeyPEM, "\n"),
			},
		}

		log.Printf("域名 %s 证书注入成功", domain)
	}

	log.Println("所有证书注入完成")
	return nil
}

// Start 启动 Xray 服务
//
// 返回:
//   - error: 如果启动失败则返回错误
func (l *Launcher) Start() error {
	if l.xrayConfig == nil {
		return fmt.Errorf("Xray 配置未加载")
	}

	log.Println("正在构建 Xray 配置...")

	// 构建 protobuf 配置
	pbConfig, err := l.xrayConfig.Build()
	if err != nil {
		return fmt.Errorf("构建 Xray 配置失败: %w", err)
	}

	log.Println("正在创建 Xray 实例...")

	// 创建 Xray 实例
	server, err := core.New(pbConfig)
	if err != nil {
		return fmt.Errorf("创建 Xray 实例失败: %w", err)
	}

	log.Println("正在启动 Xray 服务...")

	// 启动服务
	if err := server.Start(); err != nil {
		return fmt.Errorf("启动 Xray 服务失败: %w", err)
	}

	l.server = server
	log.Println("Xray 服务启动成功")

	// 触发 GC 清理配置加载产生的垃圾
	runtime.GC()
	debug.FreeOSMemory()

	return nil
}

// WaitForShutdown 等待系统信号以优雅关闭服务
// 监听 SIGINT 和 SIGTERM 信号
func (l *Launcher) WaitForShutdown() {
	log.Println("Xray 服务运行中，按 Ctrl+C 停止...")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("收到停止信号，正在关闭服务...")
}

// Close 关闭 Xray 服务并释放资源
func (l *Launcher) Close() error {
	if l.server != nil {
		log.Println("正在关闭 Xray 服务...")
		if err := l.server.Close(); err != nil {
			return fmt.Errorf("关闭 Xray 服务失败: %w", err)
		}
		log.Println("Xray 服务已关闭")
	}
	return nil
}

// Run 执行完整的启动流程
// 包括加载配置、确保证书、注入证书、启动服务并等待关闭
//
// 参数:
//   - ctx: 上下文，用于控制超时和取消
//
// 返回:
//   - error: 如果任何步骤失败则返回错误
func (l *Launcher) Run(ctx context.Context) error {
	// 1. 加载 Xray 配置
	if err := l.LoadConfig(); err != nil {
		return err
	}

	// 2. 提取域名并确保证书就绪
	domains := l.ExtractDomains()
	if err := l.certManager.EnsureCertificates(ctx, domains); err != nil {
		return err
	}

	// 3. 注入证书到配置
	if err := l.InjectCertificates(); err != nil {
		return err
	}

	// 4. 启动 Xray 服务
	if err := l.Start(); err != nil {
		return err
	}

	// 5. 等待关闭信号
	l.WaitForShutdown()

	// 6. 关闭服务
	return l.Close()
}
