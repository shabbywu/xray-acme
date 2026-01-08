// Package acme 提供 TLS 证书的管理功能
// 包括证书的自动申请、更新和提取
package acme

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"strings"
)

// CertKeyPair 证书和私钥的 PEM 格式对
type CertKeyPair struct {
	// CertPEM 证书链的 PEM 格式内容
	CertPEM string

	// KeyPEM 私钥的 PEM 格式内容
	KeyPEM string
}

// 证书相关错误
var (
	// ErrNilCertificate 证书为空错误
	ErrNilCertificate = errors.New("certificate is nil")

	// ErrNilPrivateKey 私钥为空错误
	ErrNilPrivateKey = errors.New("private key is nil")
)

// ExtractCertKeyPair 从 TLS 证书中提取证书链和私钥的 PEM 格式
//
// 参数:
//   - tlsCert: TLS 证书对象
//
// 返回:
//   - *CertKeyPair: 证书和私钥的 PEM 格式对
//   - error: 如果提取失败则返回错误
func ExtractCertKeyPair(tlsCert *tls.Certificate) (*CertKeyPair, error) {
	if tlsCert == nil {
		return nil, ErrNilCertificate
	}

	pair := &CertKeyPair{}

	// 提取证书链
	certPEM, err := extractCertificateChain(tlsCert)
	if err != nil {
		return nil, err
	}
	pair.CertPEM = certPEM

	// 提取私钥
	keyPEM, err := extractPrivateKey(tlsCert)
	if err != nil {
		return nil, err
	}
	pair.KeyPEM = keyPEM

	return pair, nil
}

// extractCertificateChain 从 TLS 证书中提取证书链的 PEM 格式
func extractCertificateChain(tlsCert *tls.Certificate) (string, error) {
	if len(tlsCert.Certificate) == 0 {
		return "", errors.New("证书链为空")
	}

	var certs []string
	for i, derCert := range tlsCert.Certificate {
		block := &pem.Block{
			Type:  "CERTIFICATE",
			Bytes: derCert,
		}
		certs = append(certs, string(pem.EncodeToMemory(block)))

		// 解析并打印第一个证书（叶子证书）的信息
		if i == 0 {
			if x509Cert, err := x509.ParseCertificate(derCert); err == nil {
				log.Printf("证书主题: %s", x509Cert.Subject)
				log.Printf("证书有效期: %s - %s", x509Cert.NotBefore, x509Cert.NotAfter)
			}
		}
	}

	return strings.Join(certs, ""), nil
}

// extractPrivateKey 从 TLS 证书中提取私钥的 PEM 格式
func extractPrivateKey(tlsCert *tls.Certificate) (string, error) {
	if tlsCert.PrivateKey == nil {
		return "", ErrNilPrivateKey
	}

	var keyBytes []byte
	var keyType string

	// 根据私钥类型选择编码方式
	switch k := tlsCert.PrivateKey.(type) {
	case interface{ MarshalPKCS1PrivateKey() ([]byte, error) }:
		// RSA 私钥使用 PKCS1 格式
		var err error
		keyBytes, err = k.MarshalPKCS1PrivateKey()
		if err != nil {
			return "", fmt.Errorf("序列化 RSA 私钥失败: %w", err)
		}
		keyType = "RSA PRIVATE KEY"
	default:
		// 其他类型使用 PKCS8 格式
		var err error
		keyBytes, err = x509.MarshalPKCS8PrivateKey(tlsCert.PrivateKey)
		if err != nil {
			return "", fmt.Errorf("序列化私钥失败: %w", err)
		}
		keyType = "PRIVATE KEY"
	}

	block := &pem.Block{
		Type:  keyType,
		Bytes: keyBytes,
	}

	return string(pem.EncodeToMemory(block)), nil
}
