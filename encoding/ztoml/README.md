# atoml - TOML 数据处理库 / TOML Data Processing Library

[中文](#中文) | [English](#english)

---

## 中文

### 📖 简介

`atoml` 是一个高性能TOML数据处理库，提供TOML编码、解码、类型安全转换及JSON格式转换等功能。  
适用于配置文件处理、数据序列化、结构化数据转换等场景。

GitHub地址: [github.com/yuancore/go-zen/encoding/atoml](https://github.com/yuancore/go-zen/encoding/atoml)

### 📦 安装

```bash
go get github.com/yuancore/go-zen/encoding/atoml
```

### 🚀 快速开始

#### 编码示例
```go
package main

import (
	"fmt"
	"github.com/yuancore/go-zen/encoding/atoml"
)

type ServerConfig struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
}

func main() {
	config := ServerConfig{
		Host: "0.0.0.0",
		Port: 8080,
	}

	// 编码为TOML
	tomlBytes, err := atoml.Encode(config)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(tomlBytes))
	
	// 输出:
	// host = "0.0.0.0"
	// port = 8080
}
```

#### 解码到结构体
```go
func main() {
	tomlContent := `
	[database]
	host = "127.0.0.1"
	port = 3306
	`

	type DatabaseConfig struct {
		Host string `toml:"host"`
		Port int    `toml:"port"`
	}

	var config struct {
		Database DatabaseConfig `toml:"database"`
	}

	err := atoml.DecodeTo([]byte(tomlContent), &config)
	if err != nil {
		panic(err)
	}
	
	fmt.Println(config.Database.Host) // 输出: 127.0.0.1
	fmt.Println(config.Database.Port) // 输出: 3306
}
```

#### 解码到Map
```go
func main() {
	tomlContent := `
	[user]
	name = "Alice"
	skills = ["Go", "Rust"]
	`

	// 解码到通用map
	result, err := atoml.Decode([]byte(tomlContent))
	if err != nil {
		panic(err)
	}
	
	name := result["user"].(map[string]interface{})["name"].(string)
	fmt.Println(name) // 输出: Alice
}
```

#### 转换为JSON
```go
func main() {
	tomlContent := `
	product = "Laptop"
	price = 1299.99
	features = ["Battery", "Touchscreen"]
	`

	jsonBytes, err := atoml.ToJson([]byte(tomlContent))
	if err != nil {
		panic(err)
	}
	
	fmt.Println(string(jsonBytes))
	// 输出: {"features":["Battery","Touchscreen"],"price":1299.99,"product":"Laptop"}
}
```

### ✨ 核心特性

| 特性                | 描述                                                                 |
|---------------------|--------------------------------------------------------------------|
| **高效编码/解码**    | 基于BurntSushi/toml库的高性能实现                                |
| **类型安全**         | 支持结构体标签映射和自动类型转换                                  |
| **内存优化**         | 使用sync.Pool复用缓冲区，减少GC压力                              |
| **格式转换**         | 一键将TOML转换为标准JSON格式                                     |
| **并发安全**         | 所有导出方法线程安全                                              |

### ⚠️ 注意事项
1. `Decode`返回的map需要手动类型断言获取具体值
2. 结构体字段需使用`toml`标签指定映射关系
3. TOML数组类型会被转换为`[]interface{}`
4. 数字溢出时会返回原始字符串值
5. 建议使用DecodeTo直接解析到结构体以获得最佳性能

### 🤝 参与贡献
[贡献指南](https://github.com/yuancore/go-zen/blob/main/CONTRIBUTING.md) | [提交Issue](https://github.com/yuancore/go-zen/issues)

---

## English

### 📖 Introduction

`atoml` is a high-performance TOML processing library providing encoding, decoding, type-safe conversions and JSON transformation.  
Suitable for configuration processing, data serialization and structured data conversion.

GitHub URL: [github.com/yuancore/go-zen/encoding/atoml](https://github.com/yuancore/go-zen/encoding/atoml)

### 📦 Installation

```bash
go get github.com/yuancore/go-zen/encoding/atoml
```

### 🚀 Quick Start

#### Encoding Example
```go
type ServerConfig struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
}

func main() {
	config := ServerConfig{
		Host: "0.0.0.0",
		Port: 8080,
	}

	tomlBytes, err := atoml.Encode(config)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(tomlBytes))
}
```

#### Decoding to Struct
```go
tomlContent := `
[database]
host = "127.0.0.1"
port = 3306
`

type DatabaseConfig struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
}

var config struct {
	Database DatabaseConfig `toml:"database"`
}

err := atoml.DecodeTo([]byte(tomlContent), &config)
```

#### Convert to JSON
```go
tomlContent := `
product = "Laptop"
price = 1299.99
features = ["Battery", "Touchscreen"]
`

jsonBytes, err := atoml.ToJson([]byte(tomlContent))
```

### ✨ Key Features

| Feature             | Description                                                     |
|---------------------|-----------------------------------------------------------------|
| **High Performance**| Built on BurntSushi/toml with sync.Pool optimizations          |
| **Type Safety**     | Struct tag mapping and auto-conversion                         |
| **Memory Efficient**| Buffer reuse minimizes GC pressure                             |
| **JSON Conversion** | Convert TOML to standard JSON with single method               |
| **Concurrency Safe**| All exported methods are thread-safe                           |

### ⚠️ Important Notes
1. Manual type assertion required for `Decode` results
2. Use `toml` tags for struct field mapping
3. TOML arrays become `[]interface{}` in Go
4. Returns original string on number overflow
5. Prefer DecodeTo for better performance with known structures

### 🤝 Contributing
[Contribution Guide](https://github.com/yuancore/go-zen/blob/main/CONTRIBUTING.md) | [Open an Issue](https://github.com/yuancore/go-zen/issues)

[⬆ Back to Top](#中文)
