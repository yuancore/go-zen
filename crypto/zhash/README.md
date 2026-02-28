# ahash - 哈希计算库 / Hash Calculation Utilities

[中文](#中文) | [English](#english)

---

## 中文

### 📖 简介

`ahash` 是一个简单高效的哈希计算库，支持多种常见的哈希算法，包括 MD5、SHA1、SHA256、SHA512 和 CRC32。适用于数据校验、密码存储、文件完整性验证等场景。

GitHub地址: [github.com/yuancore/go-zen/crypto/ahash](https://github.com/yuancore/go-zen/crypto/ahash)

### 📦 安装

```bash
go get github.com/yuancore/go-zen/crypto/ahash
```

### 🚀 快速开始

#### MD5 哈希计算
```go
package main

import (
	"fmt"
	"github.com/yuancore/go-zen/crypto/ahash"
)

func main() {
	data := "需要计算哈希的数据"
	hash := ahash.MD5(data)
	fmt.Println("MD5 哈希值:", hash)
}
```

#### SHA256 哈希计算
```go
func main() {
	data := "需要计算哈希的数据"
	hash := ahash.SHA256(data)
	fmt.Println("SHA256 哈希值:", hash)
}
```

#### CRC32 校验和计算
```go
func main() {
	data := "需要计算校验和的数据"
	checksum := ahash.Crc32(data)
	fmt.Println("CRC32 校验和:", checksum)
}
```

### 🔧 高级用法

#### 批量计算哈希
```go
func main() {
	data := "批量计算哈希的数据"
	md5Hash := ahash.MD5(data)
	sha1Hash := ahash.SHA1(data)
	sha256Hash := ahash.SHA256(data)
	sha512Hash := ahash.SHA512(data)
	crc32Checksum := ahash.Crc32(data)

	fmt.Println("MD5:", md5Hash)
	fmt.Println("SHA1:", sha1Hash)
	fmt.Println("SHA256:", sha256Hash)
	fmt.Println("SHA512:", sha512Hash)
	fmt.Println("CRC32:", crc32Checksum)
}
```

### ✨ 核心特性

| 特性                | 描述                                                                 |
|---------------------|--------------------------------------------------------------------|
| **多算法支持**       | 支持 MD5、SHA1、SHA256、SHA512 和 CRC32 等多种哈希算法             |
| **简单易用**         | 提供简洁的 API，快速计算哈希值                                     |
| **高性能**           | 基于 Go 语言原生哈希库实现，性能优异                               |
| **跨平台**           | 支持所有 Go 语言支持的平台                                         |

### ⚠️ 注意事项
1. MD5 和 SHA1 已不再推荐用于密码存储等安全场景，建议使用 SHA256 或 SHA512。
2. CRC32 主要用于校验和数据完整性验证，不适用于加密场景。
3. 哈希值不可逆，无法从哈希值还原原始数据。

### 🤝 参与贡献
[贡献指南](https://github.com/yuancore/go-zen/blob/main/CONTRIBUTING.md) | [提交Issue](https://github.com/yuancore/go-zen/issues)

---

## English

### 📖 Introduction

`ahash` is a simple and efficient hash calculation library supporting multiple common hash algorithms, including MD5, SHA1, SHA256, SHA512, and CRC32. It is suitable for data verification, password storage, file integrity checks, and more.

GitHub URL: [github.com/yuancore/go-zen/crypto/ahash](https://github.com/yuancore/go-zen/crypto/ahash)

### 📦 Installation

```bash
go get github.com/yuancore/go-zen/crypto/ahash
```

### 🚀 Quick Start

#### MD5 Hash Calculation
```go
package main

import (
	"fmt"
	"github.com/yuancore/go-zen/crypto/ahash"
)

func main() {
	data := "data to hash"
	hash := ahash.MD5(data)
	fmt.Println("MD5 Hash:", hash)
}
```

#### SHA256 Hash Calculation
```go
func main() {
	data := "data to hash"
	hash := ahash.SHA256(data)
	fmt.Println("SHA256 Hash:", hash)
}
```

#### CRC32 Checksum Calculation
```go
func main() {
	data := "data to checksum"
	checksum := ahash.Crc32(data)
	fmt.Println("CRC32 Checksum:", checksum)
}
```

### 🔧 Advanced Usage

#### Batch Hash Calculation
```go
func main() {
	data := "data to hash"
	md5Hash := ahash.MD5(data)
	sha1Hash := ahash.SHA1(data)
	sha256Hash := ahash.SHA256(data)
	sha512Hash := ahash.SHA512(data)
	crc32Checksum := ahash.Crc32(data)

	fmt.Println("MD5:", md5Hash)
	fmt.Println("SHA1:", sha1Hash)
	fmt.Println("SHA256:", sha256Hash)
	fmt.Println("SHA512:", sha512Hash)
	fmt.Println("CRC32:", crc32Checksum)
}
```

### ✨ Key Features

| Feature             | Description                                                     |
|---------------------|-----------------------------------------------------------------|
| **Multi-algorithm** | Supports MD5, SHA1, SHA256, SHA512, and CRC32                   |
| **Easy to Use**     | Provides a simple API for quick hash calculation                |
| **High Performance**| Built on Go's native hash libraries for excellent performance   |
| **Cross-platform**  | Supports all platforms compatible with Go                       |

### ⚠️ Important Notes
1. MD5 and SHA1 are no longer recommended for security-sensitive scenarios like password storage. Use SHA256 or SHA512 instead.
2. CRC32 is mainly used for checksum and data integrity verification, not for encryption.
3. Hashes are irreversible; original data cannot be restored from the hash value.

### 🤝 Contributing
[Contribution Guide](https://github.com/yuancore/go-zen/blob/main/CONTRIBUTING.md) | [Open an Issue](https://github.com/yuancore/go-zen/issues)

[⬆ Back to Top](#中文)