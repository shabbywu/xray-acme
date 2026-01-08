// Package acme 提供 TLS 证书的管理功能
// 包括证书的自动申请、更新和提取
package acme

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/caddyserver/certmagic"
	"github.com/libdns/dnspod"
	"log"

	"xray-acme/config"
)

// Manager 证书管理器
// 封装了 certmagic 的功能，提供证书的申请和管理服务
type Manager struct {
	config      *certmagic.Config
	dnsProvider *dnspod.Provider
	appConfig   *config.Config
}

// NewManager 创建新的证书管理器
// 参数:
//   - cfg: 应用程序配置
//
// 返回:
//   - *Manager: 证书管理器实例
func NewManager(cfg *config.Config) *Manager {
	// 创建 DNS 提供商实例
	dnsProvider := &dnspod.Provider{
		APIToken: cfg.GetDNSPodAPIToken(),
	}

	// 配置 CertMagic 默认 ACME 设置
	certmagic.DefaultACME.Email = cfg.Cert.Email
	certmagic.DefaultACME.DNS01Solver = &certmagic.DNS01Solver{
		DNSManager: certmagic.DNSManager{
			DNSProvider: dnsProvider,
		},
	}

	return &Manager{
		config:      certmagic.NewDefault(),
		dnsProvider: dnsProvider,
		appConfig:   cfg,
	}
}

// EnsureCertificates 确保指定域名的证书已就绪
// 如果证书不存在或已过期，将自动申请新证书
//
// 参数:
//   - ctx: 上下文，用于控制超时和取消
//   - domains: 需要证书的域名列表
//
// 返回:
//   - error: 如果证书申请失败则返回错误
func (m *Manager) EnsureCertificates(ctx context.Context, domains []string) error {
	if len(domains) == 0 {
		log.Println("没有需要申请证书的域名")
		return nil
	}

	log.Printf("正在检查/申请域名证书 (DNS Provider: DNSPod)...")
	log.Printf("域名列表: %v", domains)

	// ManageSync 会阻塞直到所有证书就绪
	if err := m.config.ManageSync(ctx, domains); err != nil {
		return fmt.Errorf("证书申请失败: %w", err)
	}

	log.Println("所有证书状态正常")
	return nil
}

// GetCertificate 获取指定域名的证书
//
// 参数:
//   - domain: 域名
//
// 返回:
//   - *tls.Certificate: TLS 证书
//   - error: 如果获取失败则返回错误
func (m *Manager) GetCertificate(domain string) (*tls.Certificate, error) {
	cert, err := m.config.GetCertificate(&tls.ClientHelloInfo{
		ServerName: domain,
	})
	if err != nil {
		return nil, fmt.Errorf("获取域名 %s 的证书失败: %w", domain, err)
	}
	return cert, nil
}

// GetCertKeyPair 获取指定域名的证书和私钥 PEM 对
//
// 参数:
//   - domain: 域名
//
// 返回:
//   - *CertKeyPair: 证书和私钥的 PEM 格式对
//   - error: 如果获取或提取失败则返回错误
func (m *Manager) GetCertKeyPair(domain string) (*CertKeyPair, error) {
	cert, err := m.GetCertificate(domain)
	if err != nil {
		return nil, err
	}
	return ExtractCertKeyPair(cert)
}
