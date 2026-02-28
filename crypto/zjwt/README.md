# jwt - JWT 工具库 / JWT Utility Library

[中文](#中文-1) | [English](#english-1)

---

## 中文-1

### 📖 简介

`jwt` 是一个轻量级、高性能的 JWT（JSON Web Token）库，专为 Go 语言设计。它支持多种签名算法（如 HS256、RS256 等），提供自动令牌刷新机制和分布式系统兼容能力，帮助开发者快速实现安全的身份认证与授权功能。

GitHub 地址: [github.com/yuancore/go-zen/utils/jwt](https://github.com/yuancore/go-zen/utils/jwt)

### 🎯 核心亮点

- **多算法支持** - 内置 HS256/RS256/ES256 等 8 种主流算法
- **自动刷新** - 支持令牌过期前自动刷新，无需手动干预
- **零依赖** - 不依赖第三方库，极简设计保障高性能
- **分布式友好** - 提供黑名单机制和集群环境验证能力
- **链式编程** - 优雅的链式调用 API 设计
- **自定义 Claims** - 灵活扩展业务专属身份载荷

### 📦 安装

```bash
go get github.com/yuancore/go-zen/utils/jwt
```

### 🚀 快速开始

#### 生成 Token

```go
package main

import (
	"fmt"
	"github.com/yuancore/go-zen/utils/jwt"
	"time"
)

func main() {
	// 初始化 JWT 实例（默认使用 HS256 算法）
	token := jwt.New()
	
	// 设置 Claims
	claims := token.Claims()
	claims.Set("user_id", 123)
	claims.SetExpiresAt(time.Now().Add(2 * time.Hour)) // 2 小时过期

	// 设置密钥（HS256 需要）
	token.SetSecret("your-secret-key")

	// 生成 Token
	signedToken, err := token.Create()
	if err != nil {
		panic(err)
	}
	fmt.Println("Generated Token:", signedToken)
}
```

#### 解析 Token（自动刷新）

```go
package main

import (
	"fmt"
	"github.com/yuancore/go-zen/utils/jwt"
)

func main() {
	// 待验证 Token
	tokenString := "your.jwt.token"

	// 初始化解析器
	parser := jwt.NewParse()
	
	// 配置验证参数
	parser.SetSecret("your-secret-key") // HS256 需要
	parser.SetAutoRefresh(true)         // 开启自动刷新
	parser.SetExpireThreshold(10 * time.Minute) // 过期前10分钟自动刷新

	// 解析 Token
	claims, newToken, err := parser.Parse(tokenString)
	if err != nil {
		panic(err)
	}

	// 输出结果
	fmt.Println("原始 Claims:", claims)
	if newToken != "" {
		fmt.Println("新 Token 已生成:", newToken)
	}
}
```

#### 使用 RS256 算法

```go
package main

import (
	"fmt"
	"github.com/yuancore/go-zen/utils/jwt"
	"time"
)

func main() {
	// 初始化并指定算法
	token := jwt.New().SetAlg(jwt.RS256)
	
	// 读取 RSA 密钥
	privateKey := jwt.MustReadKey("private.pem")
	publicKey := jwt.MustReadKey("public.pem")

	// 配置密钥与 Claims
	token.SetPrivateKey(privateKey)
	token.Claims().
		Set("role", "admin").
		SetExpiresAt(time.Now().Add(24 * time.Hour))

	// 生成 Token
	signedToken, err := token.Create()
	if err != nil {
		panic(err)
	}
	
	// 验证 Token
	parser := jwt.NewParse().
		SetAlg(jwt.RS256).
		SetPublicKey(publicKey)

	if claims, _, err := parser.Parse(signedToken); err == nil {
		fmt.Println("用户角色:", claims.Get("role"))
	}
}
```

### ✨ 核心特性

| 特性                  | 描述                                                                 |
|-----------------------|----------------------------------------------------------------------|
| **全算法支持**         | HS256/HS384/HS512/RS256/RS384/RS512/ES256/ES512                     |
| **双模式解析**         | 严格模式（验证签名+过期时间） / 宽松模式（仅验证签名）                |
| **自动刷新机制**       | 根据阈值自动生成新 Token，无缝衔接业务系统                           |
| **密钥热加载**         | 支持运行时动态更新签名密钥，满足轮转需求                             |
| **多场景适配**         | 提供 Cookie/Header 自动提取、自定义校验钩子等扩展能力                 |

### ⚠️ 注意事项

1. HS256 密钥长度建议 ≥ 256 位，RS256 密钥推荐 ≥ 2048 位
2. 生产环境务必通过安全渠道存储和传输密钥
3. 自动刷新功能需客户端配合处理新 Token 的回传
4. 黑名单功能需自行实现持久化存储

---

## English-1

### 📖 Introduction

`jwt` is a lightweight, high-performance JWT (JSON Web Token) library designed for Go. It supports multiple signing algorithms (e.g., HS256, RS256), provides automatic token refresh, and offers distributed system compatibility, enabling developers to implement secure authentication and authorization quickly.

GitHub URL: [github.com/yuancore/go-zen/utils/jwt](https://github.com/yuancore/go-zen/utils/jwt)

### 🎯 Key Highlights

- **Multi-Algorithm Support** - 8 mainstream algorithms including HS256/RS256/ES256
- **Auto-Refresh** - Automatic token refresh before expiration
- **Zero Dependency** - No third-party dependencies, minimal design for high performance
- **Distributed-Friendly** - Blacklist mechanism and cluster environment validation
- **Fluent API** - Elegant chainable method design
- **Custom Claims** - Flexible extension for business-specific payloads

### 📦 Installation

```bash
go get github.com/yuancore/go-zen/utils/jwt
```

### 🚀 Quick Start

#### Generate Token

```go
package main

import (
	"fmt"
	"github.com/yuancore/go-zen/utils/jwt"
	"time"
)

func main() {
	// Initialize JWT instance (default: HS256)
	token := jwt.New()
	
	// Set Claims
	claims := token.Claims()
	claims.Set("user_id", 123)
	claims.SetExpiresAt(time.Now().Add(2 * time.Hour))

	// Set secret (required for HS256)
	token.SetSecret("your-secret-key")

	// Generate Token
	signedToken, err := token.Create()
	if err != nil {
		panic(err)
	}
	fmt.Println("Generated Token:", signedToken)
}
```

#### Parse Token (Auto-Refresh)

```go
package main

import (
	"fmt"
	"github.com/yuancore/go-zen/utils/jwt"
)

func main() {
	// Sample token
	tokenString := "your.jwt.token"

	// Initialize parser
	parser := jwt.NewParse()
	
	// Configure validation
	parser.SetSecret("your-secret-key")     // Required for HS256
	parser.SetAutoRefresh(true)             // Enable auto-refresh
	parser.SetExpireThreshold(10 * time.Minute) // Refresh 10 mins before expiry

	// Parse Token
	claims, newToken, err := parser.Parse(tokenString)
	if err != nil {
		panic(err)
	}

	// Output results
	fmt.Println("Original Claims:", claims)
	if newToken != "" {
		fmt.Println("New Token Generated:", newToken)
	}
}
```

#### Using RS256 Algorithm

```go
package main

import (
	"fmt"
	"github.com/yuancore/go-zen/utils/jwt"
	"time"
)

func main() {
	// Initialize with RS256 algorithm
	token := jwt.New().SetAlg(jwt.RS256)
	
	// Load RSA keys
	privateKey := jwt.MustReadKey("private.pem")
	publicKey := jwt.MustReadKey("public.pem")

	// Configure keys and claims
	token.SetPrivateKey(privateKey)
	token.Claims().
		Set("role", "admin").
		SetExpiresAt(time.Now().Add(24 * time.Hour))

	// Generate Token
	signedToken, err := token.Create()
	if err != nil {
		panic(err)
	}
	
	// Verify Token
	parser := jwt.NewParse().
		SetAlg(jwt.RS256).
		SetPublicKey(publicKey)

	if claims, _, err := parser.Parse(signedToken); err == nil {
		fmt.Println("User Role:", claims.Get("role"))
	}
}
```

### ✨ Core Features

| Feature                     | Description                                                                 |
|-----------------------------|-----------------------------------------------------------------------------|
| **Full Algorithm Support**  | HS256/HS384/HS512/RS256/RS384/RS512/ES256/ES512                            |
| **Dual Validation Modes**   | Strict (signature+expiry) / Loose (signature only)                         |
| **Auto-Refresh**            | Generate new tokens before expiry without interruption                     |
| **Hot Key Reload**          | Dynamically update signing keys during runtime                             |
| **Extensible Hooks**        | Custom validation hooks, Cookie/Header extraction                          |

### ⚠️ Important Notes

1. Recommended key lengths: HS256 ≥ 256 bits, RS256 ≥ 2048 bits
2. Always store and transmit secrets securely in production
3. Auto-refresh requires client-side cooperation to handle new tokens
4. Blacklist persistence needs custom implementation

[⬆ Back to Top](#中文-1)