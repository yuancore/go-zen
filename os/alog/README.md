# alog - 高性能日志管理库 / High-Performance Logging Manager

[中文](#中文) | [English](#english)

---

## 中文

### 📖 简介

`alog` 是一个高性能、灵活且易于使用的日志管理库，专为 Go 项目设计。它基于 `zap` 日志库，提供了丰富的日志功能，包括日志级别控制、日志文件轮转、多输出目标（控制台和文件）、自定义日志格式等。`alog` 旨在为开发者提供一个简洁、高效的日志解决方案，适用于从开发到生产环境的全生命周期。

GitHub 地址: [github.com/yuancore/go-zen/os/alog](https://github.com/yuancore/go-zen/os/alog)

### 📦 安装

```bash
go get github.com/yuancore/go-zen/os/alog
```

### 🚀 快速开始

#### 初始化日志实例
```go
package main

import (
	"github.com/yuancore/go-zen/os/alog"
)

func main() {
	// 初始化日志实例
	logger := alog.New("/var/log/myapp.log")

	// 设置日志级别为 Info
	logger.SetLevel("info")

	// 注册日志器
	logger.Register()

	// 记录日志
	alog.Info("应用程序已启动")
	alog.Debug("这是一个调试信息") // 由于日志级别为 Info，此条日志不会输出
	alog.Error("发生了一个错误", alog.String("error", "出错了"))
}
```

#### 设置日志输出到控制台
```go
func main() {
	// 初始化日志实例
	logger := alog.New("/var/log/myapp.log")

	// 设置日志输出到控制台
	logger.SetConsole(true)

	// 注册日志器
	logger.Register()

	// 记录日志
	alog.Info("日志已输出到控制台")
}
```

#### 自定义日志格式
```go
func main() {
	// 初始化日志实例
	logger := alog.New("/var/log/myapp.log")

	// 设置日志格式为 JSON
	logger.SetFormat("json")

	// 注册日志器
	logger.Register()

	// 记录日志
	alog.Info("日志格式已设置为 JSON")
}
```

#### 设置日志文件轮转
```go
func main() {
	// 初始化日志实例
	logger := alog.New("/var/log/myapp.log")

	// 设置日志文件最大大小为 100MB，最多保留 10 个备份文件，最长保留 30 天
	logger.SetMaxSize(100).SetMaxBackups(10).SetMaxAge(30)

	// 注册日志器
	logger.Register()

	// 记录日志
	alog.Info("日志文件轮转设置已生效")
}
```

### ✨ 核心特性

| 特性                  | 描述                                                                 |
|-----------------------|----------------------------------------------------------------------|
| **高性能日志记录**     | 基于 `zap` 日志库，提供高性能的日志记录功能                           |
| **多日志级别支持**     | 支持 Debug、Info、Warn、Error、Panic、Fatal 等多种日志级别             |
| **日志文件轮转**       | 支持日志文件按大小、时间轮转，避免日志文件过大                        |
| **多输出目标**         | 支持同时输出日志到控制台和文件，方便调试和生产环境使用                |
| **自定义日志格式**     | 支持 JSON 和 Console 两种日志格式，满足不同场景需求                   |
| **线程安全**           | 支持并发日志记录，确保多线程环境下的数据安全                         |
| **自动清理**           | 自动清理过期的日志文件，确保磁盘空间的有效利用                       |

### ⚠️ 注意事项

1. 确保日志文件路径正确，避免因路径错误导致日志无法写入。
2. 在高并发环境下使用时，注意日志记录的线程安全性。
3. 对于大文件，建议设置合理的日志文件轮转策略，以避免日志文件过大。
4. 对于长时间运行的任务，建议在任务执行中处理错误并进行日志记录。

### 🤝 参与贡献

[贡献指南](https://github.com/yuancore/go-zen/blob/main/CONTRIBUTING.md) | [提交 Issue](https://github.com/yuancore/go-zen/issues)

---

## English

### 📖 Introduction

`alog` is a high-performance, flexible, and easy-to-use logging library designed for Go projects. Built on top of the `zap` logging library, it provides rich logging features, including log level control, log file rotation, multiple output targets (console and file), and custom log formats. `alog` aims to provide developers with a simple and efficient logging solution suitable for the entire lifecycle from development to production.

GitHub URL: [github.com/yuancore/go-zen/os/alog](https://github.com/yuancore/go-zen/os/alog)

### 📦 Installation

```bash
go get github.com/yuancore/go-zen/os/alog
```

### 🚀 Quick Start

#### Initialize Logger Instance
```go
package main

import (
	"github.com/yuancore/go-zen/os/alog"
)

func main() {
	// Initialize logger instance
	logger := alog.New("/var/log/myapp.log")

	// Set log level to Info
	logger.SetLevel("info")

	// Register logger
	logger.Register()

	// Log messages
	alog.Info("Application started")
	alog.Debug("This is a debug message") // This won't be logged because the level is set to Info
	alog.Error("An error occurred", alog.String("error", "something went wrong"))
}
```

#### Enable Console Logging
```go
func main() {
	// Initialize logger instance
	logger := alog.New("/var/log/myapp.log")

	// Enable console logging
	logger.SetConsole(true)

	// Register logger
	logger.Register()

	// Log messages
	alog.Info("Logging to console is enabled")
}
```

#### Custom Log Format
```go
func main() {
	// Initialize logger instance
	logger := alog.New("/var/log/myapp.log")

	// Set log format to JSON
	logger.SetFormat("json")

	// Register logger
	logger.Register()

	// Log messages
	alog.Info("Log format is set to JSON")
}
```

#### Configure Log File Rotation
```go
func main() {
	// Initialize logger instance
	logger := alog.New("/var/log/myapp.log")

	// Set log file rotation: max size 100MB, max 10 backups, max age 30 days
	logger.SetMaxSize(100).SetMaxBackups(10).SetMaxAge(30)

	// Register logger
	logger.Register()

	// Log messages
	alog.Info("Log file rotation settings are applied")
}
```

### ✨ Key Features

| Feature               | Description                                                           |
|-----------------------|-----------------------------------------------------------------------|
| **High-Performance Logging** | Built on `zap` for high-performance logging                            |
| **Multi-Level Logging** | Supports Debug, Info, Warn, Error, Panic, and Fatal log levels         |
| **Log File Rotation**  | Supports log file rotation by size and time to prevent oversized files |
| **Multiple Output Targets** | Logs can be written to both console and file for flexibility           |
| **Custom Log Formats** | Supports JSON and Console formats for different use cases             |
| **Thread Safety**      | Ensures thread-safe logging in concurrent environments                 |
| **Auto Cleanup**       | Automatically cleans up expired log files to save disk space           |

### ⚠️ Important Notes

1. Ensure the log file path is correct to avoid logging failures.
2. In high-concurrency environments, ensure thread safety when logging.
3. For large files, configure appropriate log rotation policies to avoid oversized log files.
4. For long-running tasks, implement proper error handling and logging.

### 🤝 Contributing

[Contribution Guide](https://github.com/yuancore/go-zen/blob/main/CONTRIBUTING.md) | [Open an Issue](https://github.com/yuancore/go-zen/issues)

[⬆ Back to Top](#中文)