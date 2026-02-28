# aaes - AES加密解密库 / AES Encryption & Decryption Utilities

[中文](#中文) | [English](#english)

---

## 中文

### 📖 简介

`aaes` 是一个安全高效的AES加密解密库，支持多种加密模式和填充方案，提供可靠的数据保护能力。支持CBC/ECB/CTR/OFB/CFB等主流加密模式，集成PKCS7/Zero/Space等多种填充策略，内置密钥长度校验和IV校验机制。  
适用于敏感数据存储、安全通信、配置文件加密等需要AES加密的场景。

GitHub地址: [github.com/yuancore/go-zen/crypto/aaes](https://github.com/yuancore/go-zen/crypto/aaes)

### 📦 安装

```bash
go get github.com/yuancore/go-zen/crypto/aaes
```

### 🚀 快速开始

#### AES-CBC加密
```go
package main

import (
	"fmt"
	"github.com/yuancore/go-zen/crypto/aaes"
)

func main() {
	key := []byte("thisis32byteslongsecretkey1234") // 32字节密钥
	iv := []byte("16bytesinitvec!")                 // 16字节IV

	// 加密
	ciphertext, err := aaes.Encrypt([]byte("敏感数据"), key, iv, aaes.ModeCBC, aaes.PaddingPKCS7)
	if err != nil {
		panic(err)
	}

	// 解密
	plaintext, err := aaes.Decrypt(ciphertext, key, iv, aaes.ModeCBC, aaes.PaddingPKCS7)
	fmt.Println("解密结果:", string(plaintext))
}
```

#### AES-ECB加密
```go
func main() {
	key := []byte("24bytekeyexample12345678") // 24字节密钥
	
	data := []byte("需要加密的文本")
	
	// ECB模式无需IV
	ciphertext, _ := aaes.Encrypt(data, key, nil, aaes.ModeECB, aaes.PaddingZero)
	
	plaintext, _ := aaes.Decrypt(ciphertext, key, nil, aaes.ModeECB, aaes.PaddingZero)
	fmt.Println(string(plaintext)) // 输出: 需要加密的文本
}
```

### 🔧 高级用法

#### 自动生成IV
```go
func generateIV() []byte {
	iv := make([]byte, aes.BlockSize) // 16字节
	if _, err := rand.Read(iv); err != nil {
		panic(err)
	}
	return iv
}

func main() {
	key := []byte("16bytessecretkey")
	iv := generateIV() // 安全随机生成IV

	ciphertext, _ := aaes.Encrypt([]byte("带随机IV的数据"), key, iv, aaes.ModeCBC, aaes.PaddingPKCS7)
}
```

#### 自定义错误处理
```go
func main() {
	// 使用错误类型进行精细处理
	_, err := aaes.Encrypt([]byte("test"), []byte("shortkey"), nil, aaes.ModeCTR, aaes.PaddingNone)
	
	if errors.Is(err, aaes.ErrInvalidKeyLength) {
		fmt.Println("密钥长度错误：", err)
	} else if errors.Is(err, aaes.ErrInvalidIVLength) {
		fmt.Println("IV长度错误：", err)
	}
}
```

### ✨ 核心特性

| 特性                | 描述                                                                 |
|---------------------|--------------------------------------------------------------------|
| **多模式支持**       | CBC/ECB/CTR/OFB/CFB等主流加密模式                                   |
| **灵活填充方案**     | PKCS7/ISO10126/ANSIX923/Zero/Space/None等填充类型                   |
| **安全校验**         | 自动检测密钥长度(16/24/32字节)和IV长度                              |
| **错误分类**         | 细粒度错误类型(无效填充/不支持模式等)                               |
| **流式加密**         | CTR/OFB/CFB模式支持流式数据加密                                     |

### ⚠️ 注意事项
1. CBC/CFB/OFB模式必须使用与块大小(16字节)相同的IV
2. ECB模式不建议用于敏感数据加密
3. PKCS7填充需要存储填充长度信息
4. 使用随机IV（CTR/OFB模式）可增强安全性
5. 密钥长度必须为16/24/32字节（对应AES-128/AES-192/AES-256）

### 🤝 参与贡献
[贡献指南](https://github.com/yuancore/go-zen/blob/main/CONTRIBUTING.md) | [提交Issue](https://github.com/yuancore/go-zen/issues)

---

## English

### 📖 Introduction

`aaes` is a secure AES encryption library supporting multiple modes and padding schemes, providing reliable data protection with CBC/ECB/CTR/OFB/CFB modes and PKCS7/Zero/Space padding strategies.

GitHub URL: [github.com/yuancore/go-zen/crypto/aaes](https://github.com/yuancore/go-zen/crypto/aaes)

### 📦 Installation

```bash
go get github.com/yuancore/go-zen/crypto/aaes
```

### 🚀 Quick Start

#### AES-CBC Encryption
```go
key := []byte("32byteslongencryptionkey1234")
iv := []byte("16bytesinitvec!")

ciphertext, _ := aaes.Encrypt([]byte("sensitive data"), key, iv, aaes.ModeCBC, aaes.PaddingPKCS7)
plaintext, _ := aaes.Decrypt(ciphertext, key, iv, aaes.ModeCBC, aaes.PaddingPKCS7)
```

#### AES-ECB Encryption
```go
key := []byte("24bytekeyexample12345678")
data := []byte("plain text")

ciphertext, _ := aaes.Encrypt(data, key, nil, aaes.ModeECB, aaes.PaddingZero)
plaintext, _ := aaes.Decrypt(ciphertext, key, nil, aaes.ModeECB, aaes.PaddingZero)
```

### 🔧 Advanced Usage

#### Random IV Generation
```go
func secureIV() []byte {
	iv := make([]byte, aes.BlockSize)
	rand.Read(iv)
	return iv
}
```

#### Error Handling
```go
_, err := aaes.Encrypt(data, invalidKey, iv, aaes.ModeCTR, aaes.PaddingNone)
switch {
case errors.Is(err, aaes.ErrInvalidKeyLength):
	// Handle key error
case errors.Is(err, aaes.ErrInvalidPadding):
	// Handle padding error
}
```

### ✨ Key Features

| Feature             | Description                                                     |
|---------------------|-----------------------------------------------------------------|
| **Multi-mode**      | CBC/ECB/CTR/OFB/CFB modes support                               |
| **Padding Schemes** | PKCS7/ISO10126/ANSIX923/Zero/Space/None                         |
| **Security Checks** | Auto key(16/24/32 bytes) & IV validation                        |
| **Error Types**     | Detailed error classification                                   |
| **Stream Support**  | CTR/OFB/CFB support stream processing                          |

### ⚠️ Important Notes
1. IV must match block size(16 bytes) for CBC/CFB/OFB modes
2. ECB not recommended for sensitive data
3. Store IV securely for CBC mode decryption
4. Use random IVs (CTR/OFB) for enhanced security
5. Key length must be 16/24/32 bytes (AES-128/192/256)

### 🤝 Contributing
[Contribution Guide](https://github.com/yuancore/go-zen/blob/main/CONTRIBUTING.md) | [Open an Issue](https://github.com/yuancore/go-zen/issues)

[⬆ Back to Top](#中文)