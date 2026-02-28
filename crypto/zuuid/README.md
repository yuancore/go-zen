# auuid - UUID 生成与操作库 / UUID Generation and Manipulation Utilities

[中文](#中文) | [English](#english)

---

## 中文

### 📖 简介

`auuid` 是一个基于 Go 语言的 UUID 生成与操作库，支持多种 UUID 版本的生成、解析和操作。适用于分布式系统、唯一标识符生成、日志追踪等场景。

GitHub地址: [github.com/yuancore/go-zen/crypto/auuid](https://github.com/yuancore/go-zen/crypto/auuid)

### 📦 安装

```bash
go get github.com/yuancore/go-zen/crypto/auuid
```

### 🚀 快速开始

#### 生成随机 UUID (版本4)
```go
package main

import (
	"fmt"
	"github.com/yuancore/go-zen/crypto/auuid"
)

func main() {
	// 生成随机 UUID
	uuid := auuid.New()
	fmt.Println("生成的 UUID:", uuid)
}
```

#### 生成时间序列 UUID (版本1)
```go
func main() {
	// 生成时间序列 UUID
	uuid, err := auuid.Create()
	if err != nil {
		fmt.Println("生成失败:", err)
		return
	}
	fmt.Println("时间序列 UUID:", uuid)
}
```

#### 批量生成 UUID
```go
func main() {
	// 批量生成 10 个 UUID
	uuids, err := auuid.BatchGenerate(10)
	if err != nil {
		fmt.Println("批量生成失败:", err)
		return
	}
	fmt.Println("批量生成的 UUID:", uuids)
}
```

### 🔧 高级用法

#### 自定义节点 ID 生成 UUID
```go
func main() {
	// 自定义节点 ID
	node := []byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66}
	uuid, err := auuid.CreateWithNode(node)
	if err != nil {
		fmt.Println("生成失败:", err)
		return
	}
	fmt.Println("自定义节点 UUID:", uuid)
}
```

#### 解析 UUID 字符串
```go
func main() {
	// 解析 UUID 字符串
	uuidStr := "f47ac10b-58cc-0372-8567-0e02b2c3d479"
	uuid, err := auuid.StringToUUID(uuidStr)
	if err != nil {
		fmt.Println("解析失败:", err)
		return
	}
	fmt.Println("解析后的 UUID:", uuid)
}
```

### ✨ 核心特性

| 特性                | 描述                                                                 |
|---------------------|--------------------------------------------------------------------|
| **多版本支持**       | 支持生成版本 1、3、4、5 的 UUID                                    |
| **高性能**           | 基于 Go 语言原生 UUID 库实现，性能优异                             |
| **批量生成**         | 支持高效批量生成 UUID，适用于高并发场景                           |
| **自定义节点**       | 支持自定义节点 ID 生成版本 1 UUID                                 |
| **跨平台**           | 支持所有 Go 语言支持的平台                                         |

### ⚠️ 注意事项
1. 版本 1 UUID 依赖于系统时间，确保系统时钟同步。
2. 批量生成时，建议根据实际需求调整并发量。
3. 自定义节点 ID 需确保全局唯一性。

### 🤝 参与贡献
[贡献指南](https://github.com/yuancore/go-zen/blob/main/CONTRIBUTING.md) | [提交Issue](https://github.com/yuancore/go-zen/issues)

---

## English

### 📖 Introduction

`auuid` is a Go-based UUID generation and manipulation library that supports multiple UUID versions. It is suitable for distributed systems, unique identifier generation, log tracing, and more.

GitHub URL: [github.com/yuancore/go-zen/crypto/auuid](https://github.com/yuancore/go-zen/crypto/auuid)

### 📦 Installation

```bash
go get github.com/yuancore/go-zen/crypto/auuid
```

### 🚀 Quick Start

#### Generate Random UUID (Version 4)
```go
package main

import (
	"fmt"
	"github.com/yuancore/go-zen/crypto/auuid"
)

func main() {
	// Generate random UUID
	uuid := auuid.New()
	fmt.Println("Generated UUID:", uuid)
}
```

#### Generate Time-Based UUID (Version 1)
```go
func main() {
	// Generate time-based UUID
	uuid, err := auuid.Create()
	if err != nil {
		fmt.Println("Generation failed:", err)
		return
	}
	fmt.Println("Time-based UUID:", uuid)
}
```

#### Batch Generate UUIDs
```go
func main() {
	// Batch generate 10 UUIDs
	uuids, err := auuid.BatchGenerate(10)
	if err != nil {
		fmt.Println("Batch generation failed:", err)
		return
	}
	fmt.Println("Batch generated UUIDs:", uuids)
}
```

### 🔧 Advanced Usage

#### Generate UUID with Custom Node ID
```go
func main() {
	// Custom node ID
	node := []byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66}
	uuid, err := auuid.CreateWithNode(node)
	if err != nil {
		fmt.Println("Generation failed:", err)
		return
	}
	fmt.Println("Custom node UUID:", uuid)
}
```

#### Parse UUID String
```go
func main() {
	// Parse UUID string
	uuidStr := "f47ac10b-58cc-0372-8567-0e02b2c3d479"
	uuid, err := auuid.StringToUUID(uuidStr)
	if err != nil {
		fmt.Println("Parsing failed:", err)
		return
	}
	fmt.Println("Parsed UUID:", uuid)
}
```

### ✨ Key Features

| Feature             | Description                                                     |
|---------------------|-----------------------------------------------------------------|
| **Multi-version**   | Supports generating UUIDs of versions 1, 3, 4, and 5           |
| **High Performance**| Built on Go's native UUID libraries for excellent performance   |
| **Batch Generation**| Efficient batch generation for high-concurrency scenarios       |
| **Custom Node**     | Supports custom node ID for version 1 UUID generation           |
| **Cross-platform**  | Supports all platforms compatible with Go                       |

### ⚠️ Important Notes
1. Version 1 UUIDs rely on system time; ensure clock synchronization.
2. Adjust concurrency levels for batch generation based on actual needs.
3. Ensure global uniqueness for custom node IDs.

### 🤝 Contributing
[Contribution Guide](https://github.com/yuancore/go-zen/blob/main/CONTRIBUTING.md) | [Open an Issue](https://github.com/yuancore/go-zen/issues)

[⬆ Back to Top](#中文)
