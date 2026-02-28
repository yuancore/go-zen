# pool - Goroutine 池管理库 / Goroutine Pool Management Library

[中文](#中文) | [English](#english)

---

## 中文

### 📖 简介

`pool` 模块是一个基于 [ants](https://github.com/panjf2000/ants) 的 Goroutine 池管理库，专为 Go 项目设计。它提供了高效的 Goroutine 池管理，支持创建、获取和管理 Goroutine 池实例，帮助开发者更好地管理并发任务，减少资源消耗并提高应用性能。

GitHub 地址: [github.com/yuancore/go-zen/utils/pool](https://github.com/yuancore/go-zen/utils/pool)

### 📦 安装

```bash
go get github.com/yuancore/go-zen/utils/pool
```

### 🚀 快速开始

#### 初始化和获取 Goroutine 池

```go
package main

import (
	"fmt"
	"github.com/yuancore/go-zen/utils/pool"
)

func main() {
	// 初始化一个 Goroutine 池，设置池大小为 5，最大任务队列大小为 50
	err := pool.New(5, 50)
	if err != nil {
		panic(err)
	}

	// 获取 Goroutine 池实例
	poolInstance := pool.JobPool

	// 使用 Goroutine 池提交任务
	err = poolInstance.Submit(func() {
		fmt.Println("任务开始执行")
	})
	if err != nil {
		panic(err)
	}

	// 关闭池
	defer poolInstance.Release()
}
```

#### 并发执行任务

```go
package main

import (
	"fmt"
	"github.com/yuancore/go-zen/utils/pool"
	"time"
)

func main() {
	// 初始化一个 Goroutine 池，设置池大小为 10，最大任务队列大小为 100
	err := pool.New(10, 100)
	if err != nil {
		panic(err)
	}

	// 获取 Goroutine 池实例
	poolInstance := pool.JobPool

	// 提交多个任务
	for i := 0; i < 5; i++ {
		err := poolInstance.Submit(func() {
			time.Sleep(1 * time.Second)
			fmt.Println("任务执行完毕")
		})
		if err != nil {
			panic(err)
		}
	}

	// 等待任务执行完毕
	defer poolInstance.Release()
}
```

#### 使用全局辅助函数设置和获取池实例

```go
package main

import (
	"fmt"
	"github.com/yuancore/go-zen/utils/pool"
)

func main() {
	// 初始化 Goroutine 池，池大小为 5，最大任务队列大小为 20
	err := pool.New(5, 20)
	if err != nil {
		panic(err)
	}

	// 获取全局池实例
	poolInstance := pool.JobPool

	// 提交任务
	err = poolInstance.Submit(func() {
		fmt.Println("任务在池中执行")
	})
	if err != nil {
		panic(err)
	}

	// 关闭池
	defer poolInstance.Release()
}
```

### ✨ 核心特性

| 特性                  | 描述                                                                 |
|-----------------------|----------------------------------------------------------------------|
| **高效的 Goroutine 池** | 基于 [ants](https://github.com/panjf2000/ants) 实现，高效管理 Goroutine 实例 |
| **并发任务提交**       | 提供简洁的接口提交并发任务，自动管理任务调度与资源回收                 |
| **线程安全**           | 支持并发访问池并提交任务，适用于高并发场景                             |
| **池管理**             | 提供池大小设置、任务队列大小设置等，支持池的动态扩展                   |
| **易于使用**           | 通过全局辅助函数简化池实例的获取与使用                               |

### ⚠️ 注意事项

1. 在高并发场景下，确保池大小和队列大小的合理配置，避免资源耗尽。
2. 使用 `Submit` 提交任务时，务必确保任务可以成功执行，以免任务阻塞池资源。
3. 在应用程序结束时，记得调用 `Release` 释放 Goroutine 池资源。

### 🤝 参与贡献

[贡献指南](https://github.com/yuancore/go-zen/blob/main/CONTRIBUTING.md) | [提交 Issue](https://github.com/yuancore/go-zen/issues)

---

## English

### 📖 Introduction

The `pool` package is a Goroutine pool management library based on [ants](https://github.com/panjf2000/ants), designed for Go projects. It provides efficient management of Goroutine pools, supporting the creation, retrieval, and management of Goroutine pool instances. This package helps developers better manage concurrent tasks, reduce resource consumption, and improve application performance.

GitHub URL: [github.com/yuancore/go-zen/utils/pool](https://github.com/yuancore/go-zen/utils/pool)

### 📦 Installation

```bash
go get github.com/yuancore/go-zen/utils/pool
```

### 🚀 Quick Start

#### Initializing and Retrieving the Goroutine Pool

```go
package main

import (
	"fmt"
	"github.com/yuancore/go-zen/utils/pool"
)

func main() {
	// Initialize a Goroutine pool with a size of 5 and a max task queue size of 50
	err := pool.New(5, 50)
	if err != nil {
		panic(err)
	}

	// Retrieve the Goroutine pool instance
	poolInstance := pool.JobPool

	// Submit a task to the Goroutine pool
	err = poolInstance.Submit(func() {
		fmt.Println("Task started executing")
	})
	if err != nil {
		panic(err)
	}

	// Release the pool when done
	defer poolInstance.Release()
}
```

#### Concurrent Task Execution

```go
package main

import (
	"fmt"
	"github.com/yuancore/go-zen/utils/pool"
	"time"
)

func main() {
	// Initialize a Goroutine pool with a size of 10 and a max task queue size of 100
	err := pool.New(10, 100)
	if err != nil {
		panic(err)
	}

	// Retrieve the Goroutine pool instance
	poolInstance := pool.JobPool

	// Submit multiple tasks
	for i := 0; i < 5; i++ {
		err := poolInstance.Submit(func() {
			time.Sleep(1 * time.Second)
			fmt.Println("Task completed")
		})
		if err != nil {
			panic(err)
		}
	}

	// Wait for tasks to finish
	defer poolInstance.Release()
}
```

#### Using Global Helper Functions to Set and Retrieve Pool Instances

```go
package main

import (
	"fmt"
	"github.com/yuancore/go-zen/utils/pool"
)

func main() {
	// Initialize the Goroutine pool with a size of 5 and a max task queue size of 20
	err := pool.New(5, 20)
	if err != nil {
		panic(err)
	}

	// Retrieve the global Goroutine pool instance
	poolInstance := pool.JobPool

	// Submit a task
	err = poolInstance.Submit(func() {
		fmt.Println("Task is executed in the pool")
	})
	if err != nil {
		panic(err)
	}

	// Release the pool when done
	defer poolInstance.Release()
}
```

### ✨ Key Features

| Feature                     | Description                                                                |
|-----------------------------|----------------------------------------------------------------------------|
| **Efficient Goroutine Pool** | Based on [ants](https://github.com/panjf2000/ants), efficiently manages Goroutine instances |
| **Concurrent Task Submission** | Provides a simple interface to submit concurrent tasks, automatically manages task scheduling and resource recycling |
| **Thread Safety**           | Supports concurrent access to the pool and task submission, suitable for high-concurrency scenarios |
| **Pool Management**         | Offers pool size and task queue size configuration, supports dynamic pool expansion |
| **Ease of Use**             | Simplifies pool instance retrieval and usage through global helper functions |

### ⚠️ Important Notes

1. In high-concurrency scenarios, ensure the pool size and queue size are configured properly to avoid resource exhaustion.
2. Ensure tasks can execute successfully when using `Submit` to avoid blocking the pool resources.
3. Remember to call `Release` to release Goroutine pool resources when the application finishes.

### 🤝 Contributing

[Contribution Guide](https://github.com/yuancore/go-zen/blob/main/CONTRIBUTING.md) | [Open an Issue](https://github.com/yuancore/go-zen/issues)

[⬆ Back to Top](#中文)