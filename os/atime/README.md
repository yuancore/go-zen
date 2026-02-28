# atime - 高效时间处理库 / Efficient Time Handling Library

[中文](#中文) | [English](#english)

---

## 中文

### 📖 简介

`atime` 是一个专为 Go 语言设计的高效时间处理库，提供丰富的时间操作方法，包括时区转换、时间戳处理、格式化输出、时间运算等。它基于 Go 原生 `time` 包进行扩展，旨在简化复杂的时间操作，提升开发效率。无论是处理国际化时区，还是进行高精度时间计算，`atime` 都能以简洁的 API 满足需求。

GitHub 地址: [github.com/yuancore/go-zen/os/atime](https://github.com/yuancore/go-zen/os/atime)

### 📦 安装

```bash
go get github.com/yuancore/go-zen/os/atime
```

### 🚀 快速开始

#### 初始化时间对象
```go
package main

import (
	"fmt"
	"github.com/yuancore/go-zen/os/atime"
	"time"
)

func main() {
	// 从当前时间创建
	t := atime.Now()
	fmt.Println(t.String()) // 输出: 2023-10-05 15:30:45

	// 从字符串解析时间
	t2 := atime.StrToTime("2023-10-05 15:30:45")
	fmt.Println(t2.UTC().String()) // 输出 UTC 时间

	// 从时间戳创建
	t3 := atime.NewFromTimeStamp(1696527045)
	fmt.Println(t3.Format("yyyy-MM-dd HH:mm:ss")) // 输出: 2023-10-05 15:30:45
}
```

#### 时区与时间戳
```go
func main() {
	t := atime.Now()

	// 转换为 UTC 时区
	utcTime := t.UTC()
	fmt.Println(utcTime.String())

	// 获取毫秒级时间戳
	fmt.Println(t.Millisecond()) // 输出: 1696527045123
}
```

#### 时间运算
```go
func main() {
	t := atime.Now()

	// 增加 1 小时
	future := t.Add(time.Hour)
	fmt.Println(future.String())

	// 计算时间差
	t1 := atime.StrToTime("2023-10-05 12:00:00")
	t2 := atime.StrToTime("2023-10-05 15:30:45")
	duration := t2.Sub(t1)
	fmt.Println(duration) // 输出: 3h30m45s
}
```

#### 高级格式化
```go
func main() {
	t := atime.Now()

	// 自定义格式化（支持中文周显示）
	fmt.Println(t.Format("yyyy年MM月dd日 E", true)) // 输出: 2023年10月05日 星期四

	// 毫秒精度格式化
	fmt.Println(t.Format("yyyy-MM-dd HH:mm:ss.SSS")) // 输出: 2023-10-05 15:30:45.123
}
```

### ✨ 核心特性

| 特性                  | 描述                                                                 |
|-----------------------|----------------------------------------------------------------------|
| **多时区无缝转换**     | 支持 UTC 和本地时区转换，满足国际化需求                              |
| **时间戳灵活处理**     | 提供秒、毫秒、微秒、纳秒级时间戳获取                                 |
| **链式时间运算**       | 支持 `Add`、`Sub`、`Truncate` 等方法，轻松实现时间增减与截断         |
| **智能格式化**         | 支持类似 `yyyy-MM-dd HH:mm:ss` 的易读格式，兼容中英文周显示          |
| **时间范围操作**       | 快速获取某时刻的起始与结束时间（如一天、一周、一月的开始/结束）       |
| **高性能底层实现**     | 基于原生 `time` 包优化，零额外内存分配                               |

### ⚠️ 注意事项

1. **时区敏感操作**：跨时区操作时建议显式调用 `UTC()` 或 `Local()` 方法。
2. **格式化字符**：格式化字符串需使用特定占位符（如 `yyyy` 代表年份）。
3. **时间解析**：字符串解析需严格匹配 `2006-01-02 15:04:05` 格式。
4. **并发安全**：时间对象非线程安全，高并发场景建议使用独立实例。

### 🤝 参与贡献

[贡献指南](https://github.com/yuancore/go-zen/blob/main/CONTRIBUTING.md) | [提交 Issue](https://github.com/yuancore/go-zen/issues)

---

## English

### 📖 Introduction

`atime` is a high-efficiency time handling library for Go, offering comprehensive time operations such as timezone conversion, timestamp processing, formatting, and arithmetic. Built as an extension of Go's native `time` package, it simplifies complex time manipulations with an intuitive API. Whether dealing with international timezones or high-precision calculations, `atime` delivers efficiency and clarity.

GitHub URL: [github.com/yuancore/go-zen/os/atime](https://github.com/yuancore/go-zen/os/atime)

### 📦 Installation

```bash
go get github.com/yuancore/go-zen/os/atime
```

### 🚀 Quick Start

#### Initialize Time Object
```go
package main

import (
	"fmt"
	"github.com/yuancore/go-zen/os/atime"
)

func main() {
	// Create from current time
	t := atime.Now()
	fmt.Println(t.String()) // Output: 2023-10-05 15:30:45

	// Parse from string
	t2 := atime.StrToTime("2023-10-05 15:30:45")
	fmt.Println(t2.UTC().String())

	// Create from timestamp
	t3 := atime.NewFromTimeStamp(1696527045)
	fmt.Println(t3.Format("yyyy-MM-dd HH:mm:ss")) // Output: 2023-10-05 15:30:45
}
```

#### Timezone & Timestamp
```go
func main() {
	t := atime.Now()

	// Convert to UTC
	utcTime := t.UTC()
	fmt.Println(utcTime.String())

	// Get millisecond timestamp
	fmt.Println(t.Millisecond()) // Output: 1696527045123
}
```

#### Time Arithmetic
```go
func main() {
	t := atime.Now()

	// Add 1 hour
	future := t.Add(time.Hour)
	fmt.Println(future.String())

	// Calculate duration between two times
	t1 := atime.StrToTime("2023-10-05 12:00:00")
	t2 := atime.StrToTime("2023-10-05 15:30:45")
	duration := t2.Sub(t1)
	fmt.Println(duration) // Output: 3h30m45s
}
```

#### Advanced Formatting
```go
func main() {
	t := atime.Now()

	// Custom format with Chinese weekday
	fmt.Println(t.Format("yyyy-MM-dd E", true)) // Output: 2023-10-05 星期四

	// Millisecond precision
	fmt.Println(t.Format("yyyy-MM-dd HH:mm:ss.SSS")) // Output: 2023-10-05 15:30:45.123
}
```

### ✨ Key Features

| Feature               | Description                                                           |
|-----------------------|-----------------------------------------------------------------------|
| **Seamless Timezone Conversion** | Convert between UTC and local timezones effortlessly                 |
| **Timestamp Flexibility**         | Get timestamps in seconds, milliseconds, microseconds, or nanoseconds|
| **Chained Time Operations**       | Methods like `Add`, `Sub`, `Truncate` for intuitive manipulations    |
| **Smart Formatting**              | Human-readable formats (e.g., `yyyy-MM-dd`) with multilingual support|
| **Time Range Utilities**          | Get start/end of day, week, month, etc., in one line                 |
| **High-Performance Core**         | Optimized on Go's native `time` package with zero extra allocations  |

### ⚠️ Important Notes

1. **Timezone Awareness**: Explicitly use `UTC()` or `Local()` for cross-timezone operations.
2. **Format Placeholders**: Use specific tokens like `yyyy` for year in format strings.
3. **Time Parsing**: Input strings must strictly match `2006-01-02 15:04:05` format.
4. **Concurrency**: Time objects are not thread-safe; use separate instances in concurrent code.

### 🤝 Contributing

[Contribution Guide](https://github.com/yuancore/go-zen/blob/main/CONTRIBUTING.md) | [Open an Issue](https://github.com/yuancore/go-zen/issues)

[⬆ Back to Top](#中文)