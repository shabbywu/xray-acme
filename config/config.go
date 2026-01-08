// Package config 提供应用程序的配置管理功能
// 负责从环境变量、配置文件等来源加载和验证配置
package config

import (
	"errors"
	"os"
)

// Config 应用程序的主配置结构体
// 包含 DNS 提供商配置、证书管理配置和 Xray 配置
type Config struct {
	// DNS 提供商配置
	DNS DNSConfig

	// 证书管理配置
	Cert CertConfig

	// Xray 配置
	Xray XrayConfig
}

// DNSConfig DNS 提供商配置
// 当前支持 DNSPod 作为 DNS-01 挑战的 DNS 提供商
type DNSConfig struct {
	// ProviderID DNSPod 的 API ID
	ProviderID string

	// ProviderToken DNSPod 的 API Token
	ProviderToken string
}

// CertConfig 证书管理配置
type CertConfig struct {
	// Email Let's Encrypt 注册邮箱
	// 用于接收证书过期通知等重要信息
	Email string
}

// XrayConfig Xray 相关配置
type XrayConfig struct {
	// ConfigPath Xray 配置文件路径
	ConfigPath string
}

// 环境变量名称常量
const (
	// EnvDNSPodID DNSPod API ID 环境变量名
	EnvDNSPodID = "DP_ID"

	// EnvDNSPodToken DNSPod API Token 环境变量名
	EnvDNSPodToken = "DP_TOKEN"

	// EnvEmail Let's Encrypt 注册邮箱环境变量名
	EnvEmail = "EMAIL"

	// EnvXrayConfig Xray 配置文件路径环境变量名（可选）
	EnvXrayConfig = "XRAY_CONFIG"

	// EnvCertStorage 证书存储路径环境变量名（可选）
	EnvCertStorage = "CERT_STORAGE"
)

// 默认配置值
const (
	// DefaultXrayConfigPath 默认 Xray 配置文件路径
	DefaultXrayConfigPath = "config.json"
)

// 配置验证错误
var (
	// ErrMissingDNSPodID DNSPod ID 未配置错误
	ErrMissingDNSPodID = errors.New("missing required environment variable: " + EnvDNSPodID)

	// ErrMissingDNSPodToken DNSPod Token 未配置错误
	ErrMissingDNSPodToken = errors.New("missing required environment variable: " + EnvDNSPodToken)

	// ErrMissingEmail Email 未配置错误
	ErrMissingEmail = errors.New("missing required environment variable: " + EnvEmail)
)

// LoadFromEnv 从环境变量加载配置
// 返回加载的配置和可能的错误
func LoadFromEnv() (*Config, error) {
	cfg := &Config{
		DNS: DNSConfig{
			ProviderID:    os.Getenv(EnvDNSPodID),
			ProviderToken: os.Getenv(EnvDNSPodToken),
		},
		Cert: CertConfig{
			Email: os.Getenv(EnvEmail),
		},
		Xray: XrayConfig{
			ConfigPath: getEnvWithDefault(EnvXrayConfig, DefaultXrayConfigPath),
		},
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate 验证配置的完整性和有效性
// 检查所有必需的配置项是否已正确设置
func (c *Config) Validate() error {
	if c.DNS.ProviderID == "" {
		return ErrMissingDNSPodID
	}
	if c.DNS.ProviderToken == "" {
		return ErrMissingDNSPodToken
	}
	if c.Cert.Email == "" {
		return ErrMissingEmail
	}
	return nil
}

// GetDNSPodAPIToken 返回格式化的 DNSPod API Token
// DNSPod API 要求 Token 格式为 "ID,Token"
func (c *Config) GetDNSPodAPIToken() string {
	return c.DNS.ProviderID + "," + c.DNS.ProviderToken
}

// getEnvWithDefault 获取环境变量，如果不存在则返回默认值
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
