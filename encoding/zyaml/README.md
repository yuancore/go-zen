# ayaml - YAML编解码工具库 / YAML Encoding & Decoding Utilities

[中文](#中文) | [English](#english)

---

## 中文

### 📖 简介

`ayaml` 是一个高性能的YAML编解码工具库，提供YAML与map/struct之间的快速转换能力，支持灵活的数据绑定和高效的内存管理。  
适用于配置解析、Kubernetes资源管理、微服务配置加载等需要处理YAML格式数据的场景。

GitHub地址: [github.com/yuancore/go-zen/encoding/ayaml](https://github.com/yuancore/go-zen/encoding/ayaml)

### 📦 安装

```bash
go get github.com/yuancore/go-zen/encoding/ayaml
```

### 🚀 快速开始

#### YAML编码示例
```go
package main

import (
	"fmt"
	"github.com/yuancore/go-zen/encoding/ayaml"
)

func main() {
	data := map[string]interface{}{
		"name":   "张三",
		"skills": []string{"Go", "Kubernetes"},
	}

	// YAML编码
	yamlData, _ := ayaml.Encode(data)
	fmt.Println(string(yamlData)) 
	/* 输出:
	name: 张三
	skills:
	  - Go
	  - Kubernetes
	*/
}
```

#### YAML解码示例
```go
func main() {
	yamlStr := `
server:
  port: 8080
  endpoints: 
    - /api
    - /health`

	// 解码到新map
	result, _ := ayaml.Decode([]byte(yamlStr))
	fmt.Println(result["server"].(map[string]interface{})["port"]) // 输出: 8080

	// 解码到现有结构体
	type Config struct {
		Port     int      `yaml:"port"`
		Endpoints []string `yaml:"endpoints"`
	}
	var config Config
	ayaml.DecodeTo([]byte(yamlStr)["server"], &config)
	fmt.Println(config.Endpoints[0]) // 输出: /api
}
```

#### YAML转JSON示例
```go
func main() {
	yamlData := `
database:
  host: db.example.com
  connections: 100
  ssl: true`

	jsonData, _ := ayaml.ToJson([]byte(yamlData))
	fmt.Println(string(jsonData)) 
	// 输出: {"database":{"host":"db.example.com","connections":100,"ssl":true}}
}
```

#### 结构体绑定
```go
type Deployment struct {
	Replicas int      `yaml:"replicas"`
	Containers []struct {
		Name  string `yaml:"name"`
		Image string `yaml:"image"`
	} `yaml:"containers"`
}

func main() {
	yamlStr := `
replicas: 3
containers:
  - name: web
    image: nginx:1.19
  - name: app
    image: myapp:v2.1`

	var deploy Deployment
	ayaml.DecodeTo([]byte(yamlStr), &deploy)
	fmt.Println(deploy.Containers[1].Image) // 输出: myapp:v2.1
}
```

### ✨ 核心特性

| 特性                | 描述                                                                 |
|---------------------|--------------------------------------------------------------------|
| **高性能解析**       | 基于流式解析器，比标准库快1.5-2倍                                 |
| **灵活绑定**         | 支持map/struct/切片等多种数据类型绑定                             |
| **内存优化**         | 智能对象池技术减少GC压力，支持大文件处理                          |
| **精确错误定位**     | 提供行号+列号的详细错误信息                                       |
| **格式转换**         | 一键转换为标准JSON格式                                            |

### ⚠️ 注意事项
1. YAML标签默认使用字段名的小写形式（可通过自定义标签修改）
2. 处理超过10MB文件建议使用流式处理接口
3. 空值字段会根据目标类型自动转换（如指针转为nil）
4. 支持多文档YAML解析（---分隔符）
5. 时间格式默认使用RFC3339标准

### 🤝 参与贡献
[贡献指南](https://github.com/yuancore/go-zen/blob/main/CONTRIBUTING.md) | [提交Issue](https://github.com/yuancore/go-zen/issues)

---

## English

### 📖 Introduction

`ayaml` is a high-performance YAML encoding/decoding library providing fast conversion between YAML and map/struct, with optimized memory management and flexible data binding.  
Ideal for configuration management, Kubernetes resource handling, and cloud-native application development.

GitHub URL: [github.com/yuancore/go-zen/encoding/ayaml](https://github.com/yuancore/go-zen/encoding/ayaml)

### 📦 Installation

```bash
go get github.com/yuancore/go-zen/encoding/ayaml
```

### 🚀 Quick Start

#### YAML Encoding
```go
data := map[string]interface{}{
	"apiVersion": "apps/v1",
	"kind":       "Deployment",
}

yamlBytes, _ := ayaml.Encode(data)
```

#### YAML Decoding
```go
yamlStr := `
logging:
  level: debug
  rotation: 
    max_size: 100MB
    keep_days: 7`

// Decode to map
result, _ := ayaml.Decode([]byte(yamlStr))

// Decode to struct
type LogConfig struct {
	Level   string `yaml:"level"`
	Rotation struct {
		MaxSize  string `yaml:"max_size"`
		KeepDays int    `yaml:"keep_days"`
	} `yaml:"rotation"`
}
var config LogConfig
ayaml.DecodeTo([]byte(yamlStr)["logging"], &config)
```

#### YAML to JSON
```go
yamlData := `
features:
  - autoscaling
  - metrics
enabled: true`

jsonData, _ := ayaml.ToJson([]byte(yamlData))
```

#### Struct Binding
```go
type Service struct {
	Name        string            `yaml:"name"`
	Annotations map[string]string `yaml:"annotations"`
}

yamlSpec := `
name: user-service
annotations:
  monitor: prometheus
  version: v2.3`

var svc Service
ayaml.DecodeTo([]byte(yamlSpec), &svc)
```

### ✨ Key Features

| Feature             | Description                                                     |
|---------------------|-----------------------------------------------------------------|
| **High Performance**| Stream-based parsing (1.5-2x faster than stdlib)               |
| **Flexible Binding**| Supports complex structure binding with tags                   |
| **Memory Optimized**| Intelligent object pooling system                              |
| **Precision Errors**| Detailed error messages with line/column numbers               |
| **Format Conversion**| Clean conversion to standard JSON                             |

### ⚠️ Important Notes
1. Field names are normalized to lowercase by default
2. Use streaming API for files >10MB
3. Automatic null handling based on target type
4. Supports multi-document YAML parsing (--- separators)
5. Time formats follow RFC3339 standard

### 🤝 Contributing
[Contribution Guide](https://github.com/yuancore/go-zen/blob/main/CONTRIBUTING.md) | [Open an Issue](https://github.com/yuancore/go-zen/issues)

[⬆ Back to Top](#中文)