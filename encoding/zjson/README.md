# antgo/encoding/ajson - JSON 数据处理库 / JSON Data Processing Library

[中文](#中文) | [English](#english)

---

## 中文

### 📖 简介

`antgo/encoding/ajson` 是一个高性能JSON数据处理库，提供JSON解析、编码、路径查询及类型安全转换等功能。  
适用于配置管理、API数据处理、结构化数据转换等场景。

GitHub地址: [github.com/yuancore/go-zen/encoding/ajson](https://github.com/yuancore/go-zen/encoding/ajson)

### 📦 安装

```bash
go get github.com/yuancore/go-zen/encoding/ajson
```

### 🚀 快速开始

#### 基本解析
```go
package main

import (
	"fmt"
	"github.com/yuancore/go-zen/encoding/ajson"
)

func main() {
	// 示例JSON内容
	jsonContent := `{
		"user": {
			"name": "John",
			"age": 30,
			"skills": ["Go", "Python"]
		},
		"server": {
			"host": "0.0.0.0",
			"ports": [8080, 8081]
		}
	}`

	// 解析JSON内容
	j := ajson.Decode([]byte(jsonContent))
	
	// 获取嵌套值
	name := j.Get("user.name").String()
	age := j.Get("user.age").Int()
	skill := j.Get("user.skills.0").String()
	
	fmt.Println(name)  // 输出: John
	fmt.Println(age)   // 输出: 30
	fmt.Println(skill) // 输出: Go
}
```

#### 编码生成
```go
func main() {
	// 构建数据
	data := map[string]interface{}{
		"database": map[string]interface{}{
			"host": "127.0.0.1",
			"port": 3306,
		},
	}

	// 编码为JSON字符串
	jsonStr := ajson.Encode(data)
	fmt.Println(jsonStr) // 输出: {"database":{"host":"127.0.0.1","port":3306}}
}
```

#### 类型转换
```go
func main() {
	jsonContent := `{"count": "100", "active": true}`

	j := ajson.Decode([]byte(jsonContent))
	
	// 自动类型转换
	count := j.Get("count").Int()     // 字符串转数字
	active := j.Get("active").Bool()  // 布尔值转换
	
	fmt.Println(count)  // 输出: 100
	fmt.Println(active) // 输出: true
}
```

### ✨ 核心特性

| 特性                | 描述                                                                 |
|---------------------|--------------------------------------------------------------------|
| **高效解析**         | 基于标准库的高性能解析，支持复杂数据结构                           |
| **链式操作**         | 支持 `Get("path.to.value").Int()` 链式调用                         |
| **类型安全**         | 提供20+种安全类型转换方法（String/Int/Float64/Map/Array等）        |
| **零拷贝处理**       | 解析过程最小化内存分配                                             |
| **并发安全**         | 所有操作线程安全                                                   |

### ⚠️ 注意事项
1. 使用`Get()`方法时路径不存在会返回零值
2. 类型转换失败时返回对应类型的零值（如字符串转数字失败返回0）
3. 复杂JSON建议先进行`Get()`路径验证
4. 解析失败会触发panic，生产环境建议配合recover使用

### 🤝 参与贡献
[贡献指南](https://github.com/yuancore/go-zen/blob/main/CONTRIBUTING.md) | [提交Issue](https://github.com/yuancore/go-zen/issues)

---

## English

### 📖 Introduction

`antgo/encoding/ajson` is a high-performance JSON processing library providing parsing, encoding, path querying and type-safe conversions.  
Ideal for configuration management, API data processing and structured data transformation.

GitHub URL: [github.com/yuancore/go-zen/encoding/ajson](https://github.com/yuancore/go-zen/encoding/ajson)

### 📦 Installation

```bash
go get github.com/yuancore/go-zen/encoding/ajson
```

### 🚀 Quick Start

#### Basic Parsing
```go
package main

import (
	"fmt"
	"github.com/yuancore/go-zen/encoding/ajson"
)

func main() {
	jsonContent := `{
		"user": {
			"name": "John",
			"age": 30,
			"skills": ["Go", "Python"]
		},
		"server": {
			"host": "0.0.0.0",
			"ports": [8080, 8081]
		}
	}`

	j := ajson.Decode([]byte(jsonContent))
	
	name := j.Get("user.name").String()
	age := j.Get("user.age").Int()
	skill := j.Get("user.skills.0").String()
	
	fmt.Println(name)  // Output: John
	fmt.Println(age)   // Output: 30
	fmt.Println(skill) // Output: Go
}
```

#### JSON Generation
```go
func main() {
	data := map[string]interface{}{
		"database": map[string]interface{}{
			"host": "127.0.0.1",
			"port": 3306,
		},
	}

	jsonStr := ajson.Encode(data)
	fmt.Println(jsonStr) // Output: {"database":{"host":"127.0.0.1","port":3306}}
}
```

#### Type Conversion
```go
func main() {
	jsonContent := `{"count": "100", "active": true}`

	j := ajson.Decode([]byte(jsonContent))
	
	count := j.Get("count").Int()     // String to int
	active := j.Get("active").Bool()  // Boolean conversion
	
	fmt.Println(count)  // Output: 100
	fmt.Println(active) // Output: true
}
```

### ✨ Key Features

| Feature             | Description                                                        |
|---------------------|--------------------------------------------------------------------|
| **High Performance**| Standard library-based parsing with complex data support          |
| **Chained API**     | Method chaining like `Get("path.to.value").Int()`                 |
| **Type Safety**     | 20+ type conversion methods (String/Int/Float64/Map/Array etc.)   |
| **Zero-Copy**       | Minimal memory allocation during parsing                         |
| **Concurrency Safe**| All operations are thread-safe                                    |

### ⚠️ Important Notes
1. `Get()` returns zero-value for non-existent paths
2. Type conversion failures return type's zero-value
3. Validate paths with `Get()` for complex JSON
4. Parse errors trigger panic, use with recover in production

### 🤝 Contributing
[Contribution Guide](https://github.com/yuancore/go-zen/blob/main/CONTRIBUTING.md) | [Open an Issue](https://github.com/yuancore/go-zen/issues)

[⬆ Back to Top](#中文)