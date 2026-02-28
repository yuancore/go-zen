# axml - XML编解码工具库 / XML Encoding & Decoding Utilities

[中文](#中文) | [English](#english)

---

## 中文

### 📖 简介

`axml` 是一个高性能的XML编解码工具库，提供XML与map/struct之间的快速转换能力，支持灵活的数据绑定和高效的流式处理。  
适用于配置解析、API数据交换、微服务通信等需要处理XML格式数据的场景。

GitHub地址: [github.com/yuancore/go-zen/encoding/axml](https://github.com/yuancore/go-zen/encoding/axml)

### 📦 安装

```bash
go get github.com/yuancore/go-zen/encoding/axml
```

### 🚀 快速开始

#### XML编码示例
```go
package main

import (
	"fmt"
	"github.com/yuancore/go-zen/encoding/axml"
)

func main() {
	data := map[string]string{
		"name":  "张三",
		"email": "zhangsan@example.com",
	}

	// XML编码
	xmlData, _ := axml.Encode(data)
	fmt.Println(string(xmlData)) 
	// 输出: <root><name>张三</name><email>zhangsan@example.com</email></root>
}
```

#### XML解码示例
```go
func main() {
	xmlStr := `<profile><age>30</age><city>北京</city></profile>`

	// 解码到新map
	result, _ := axml.Decode([]byte(xmlStr))
	fmt.Println(result["city"]) // 输出: 北京

	// 解码到现有map
	existingMap := make(map[string]string)
	axml.DecodeTo([]byte(xmlStr), existingMap)
	fmt.Println(existingMap["age"]) // 输出: 30
}
```

#### XML转JSON示例
```go
func main() {
	xmlData := `
	<user>
		<id>1001</id>
		<preferences>Go,Java</preferences>
	</user>`

	jsonData, _ := axml.ToJson([]byte(xmlData))
	fmt.Println(string(jsonData)) 
	// 输出: {"id":"1001","preferences":"Go,Java"}
}
```

#### 结构体绑定
```go
type User struct {
	Name  string `xml:"name"`
	Roles []string `xml:"roles>role"`
}

func main() {
	xmlStr := `
	<user>
		<name>李四</name>
		<roles>
			<role>admin</role>
			<role>developer</role>
		</roles>
	</user>`

	var user User
	axml.DecodeTo([]byte(xmlStr), &user)
	fmt.Println(user.Roles[1]) // 输出: developer
}
```

### ✨ 核心特性

| 特性                | 描述                                                                 |
|---------------------|--------------------------------------------------------------------|
| **高性能解析**       | 基于Token流式处理，比标准库快2-3倍                                 |
| **灵活绑定**         | 支持map/struct/切片等多种数据类型绑定                             |
| **内存优化**         | 对象池技术减少GC压力，支持大文件处理                              |
| **错误处理**         | 提供行号定位的详细错误信息                                        |
| **格式转换**         | 一键转换为标准JSON格式                                            |

### ⚠️ 注意事项
1. XML标签名强制转换为小写字母形式
2. 处理超1MB文件建议使用`DecodeTo`复用内存
3. 空值字段默认转换为空字符串
4. CDATA内容会保留原始格式
5. 属性解析需使用`attr`标签标注

### 🤝 参与贡献
[贡献指南](https://github.com/yuancore/go-zen/blob/main/CONTRIBUTING.md) | [提交Issue](https://github.com/yuancore/go-zen/issues)

---

## English

### 📖 Introduction

`axml` is a high-performance XML encoding/decoding library providing fast conversion between XML and map/struct, with support for flexible data binding and efficient stream processing.  
Ideal for configuration parsing, API data exchange, microservice communication and other XML processing scenarios.

GitHub URL: [github.com/yuancore/go-zen/encoding/axml](https://github.com/yuancore/go-zen/encoding/axml)

### 📦 Installation

```bash
go get github.com/yuancore/go-zen/encoding/axml
```

### 🚀 Quick Start

#### XML Encoding
```go
data := map[string]interface{}{
	"user": map[string]string{
		"name":  "John",
		"email": "john@example.com",
	},
}

xmlBytes, _ := axml.Encode(data)
```

#### XML Decoding
```go
xmlStr := `<config><debug>true</debug><timeout>30</timeout></config>`

// Decode to new map
result, _ := axml.Decode([]byte(xmlStr))

// Decode to existing struct
var config struct {
	Debug   bool `xml:"debug"`
	Timeout int  `xml:"timeout"`
}
axml.DecodeTo([]byte(xmlStr), &config)
```

#### XML to JSON
```go
xmlData := `<item><id>5001</id><inStock>true</inStock></item>`
jsonData, _ := axml.ToJson([]byte(xmlData))
```

#### Struct Binding
```go
type Order struct {
	ID    string   `xml:"id"`
	Items []string `xml:"items>item"`
}

xmlOrder := `
<order>
	<id>2001</id>
	<items>
		<item>Laptop</item>
		<item>Mouse</item>
	</items>
</order>`

var order Order
axml.DecodeTo([]byte(xmlOrder), &order)
```

### ✨ Key Features

| Feature             | Description                                                     |
|---------------------|-----------------------------------------------------------------|
| **High Performance**| Token-based streaming processing (2-3x faster than stdlib)     |
| **Flexible Binding**| Supports map/struct/slice binding                              |
| **Memory Optimized**| Object pool reduces GC pressure                                |
| **Error Handling**  | Detailed error messages with line numbers                      |
| **Format Conversion**| One-click conversion to JSON                                 |

### ⚠️ Important Notes
1. XML tags are normalized to lowercase
2. Use `DecodeTo` for files >1MB
3. Empty fields return as blank strings
4. CDATA sections preserve original formatting
5. Use `attr` tag for attribute parsing

### 🤝 Contributing
[Contribution Guide](https://github.com/yuancore/go-zen/blob/main/CONTRIBUTING.md) | [Open an Issue](https://github.com/yuancore/go-zen/issues)

[⬆ Back to Top](#中文)