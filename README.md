# Xray ACME

自动管理 TLS 证书的 Xray 启动器。使用 CertMagic 和 DNSPod DNS-01 挑战自动申请和续期 Let's Encrypt 证书。

## 功能特性

- 🔐 **自动证书管理**: 自动申请和续期 Let's Encrypt TLS 证书
- 🌐 **DNS-01 挑战**: 使用 DNSPod 作为 DNS 提供商进行域名验证
- ⚡  **零配置证书**: 从 Xray 配置中自动提取域名并申请证书
- 🔄 **证书自动注入**: 自动将证书注入到 Xray TLS 配置中
- 🛡️ **优雅关闭**: 支持 SIGINT/SIGTERM 信号优雅关闭
- 🎯 **灵活配置**: 支持命令行参数和环境变量两种配置方式

## 快速开始

### 使用方法

#### 查看帮助

```bash
./xray-acme --help
```

输出：

```
Xray ACME - 自动管理 TLS 证书的 Xray 启动器

使用 CertMagic 和 DNSPod DNS-01 挑战自动申请和续期 Let's Encrypt 证书，
并将证书自动注入到 Xray 配置中，实现全自动的 TLS 证书管理。

示例:
  # 使用命令行参数
  xray-acme --config config.json --dp-id YOUR_ID --dp-token YOUR_TOKEN --email your@email.com

  # 使用环境变量
  export DP_ID="your-id"
  export DP_TOKEN="your-token"
  export EMAIL="your@email.com"
  xray-acme

Usage:
  xray-acme [flags]

Flags:
  -c, --config string     Xray 配置文件路径 (default "config.json")
      --dp-id string      DNSPod API ID (也可通过环境变量 DP_ID 设置)
      --dp-token string   DNSPod API Token (也可通过环境变量 DP_TOKEN 设置)
      --email string      Let's Encrypt 注册邮箱 (也可通过环境变量 EMAIL 设置)
  -h, --help              help for xray-acme
  -v, --version           version for xray-acme
```

#### 使用命令行参数

```bash
./xray-acme \
  --config config.json \
  --dp-id "your-dnspod-id" \
  --dp-token "your-dnspod-token" \
  --email "your-email@example.com"
```

#### 使用环境变量

```bash
export DP_ID="your-dnspod-id"
export DP_TOKEN="your-dnspod-token"
export EMAIL="your-email@example.com"

./xray-acme --config config.json
```

#### 混合使用（命令行参数优先）

```bash
export DP_ID="default-id"
export DP_TOKEN="default-token"

# 命令行参数会覆盖环境变量
./xray-acme --dp-id "override-id" --email "your@email.com"
```

### 配置参数

| 参数 | 环境变量 | 必需 | 默认值 | 说明 |
|------|----------|------|--------|------|
| `--config, -c` | `XRAY_CONFIG` | 否 | `config.json` | Xray 配置文件路径 |
| `--dp-id` | `DP_ID` | 是 | - | DNSPod API ID |
| `--dp-token` | `DP_TOKEN` | 是 | - | DNSPod API Token |
| `--email` | `EMAIL` | 是 | - | Let's Encrypt 注册邮箱 |

### Xray 配置

完整配置详看 [Xray-core](https://github.com/XTLS/Xray-core/)

在 `config.json` 中配置入站连接，程序会自动从 `streamSettings.tlsSettings.serverName` 提取域名并申请证书：

```json
{
  "inbounds": [
    {
      "port": 443,
      "protocol": "vless",
      "streamSettings": {
        "network": "tcp",
        "security": "tls",
        "tlsSettings": {
          "serverName": "your-domain.com"
        }
      }
    }
  ]
}
```

## 获取 DNSPod API 密钥

1. 登录 [DNSPod 控制台](https://console.dnspod.cn/)
2. 进入 **账号中心** -> **密钥管理**
3. 创建新的 API 密钥
4. 记录 `ID` 和 `Token`

## 开发

### 环境要求

- Go 1.25+

### 项目结构

```
xray-acme/
├── main.go              # 程序入口
├── cmd/
│   └── root.go          # Cobra 根命令定义
├── config/
│   └── config.go        # 配置管理模块
├── cert/
│   └── cert.go          # 证书管理模块
├── launcher/
│   └── launcher.go      # Xray 启动器模块
├── config.json            # Xray 配置文件示例
├── go.mod
├── go.sum
└── README.md
```

### 构建

```bash
go build -o xray-acme .
```

### 测试

```bash
go test ./...
```

### 代码格式化

```bash
go fmt ./...
```

## 许可证

MPL-2.0 license