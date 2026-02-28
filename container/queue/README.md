# queue - 异步任务队列库 / Asynchronous Task Queue Utilities

[中文](#中文) | [English](#english)

---

## 中文

### 📖 简介

`queue` 是基于 Redis 和 Go 的高性能异步任务队列库，支持延迟任务、唯一性任务、优先级队列等特性。通过 `asynq` 实现底层队列管理，提供线程安全的客户端和服务端操作，适用于分布式系统任务调度。

GitHub地址: [github.com/yuancore/go-zen/container/queue](https://github.com/yuancore/go-zen/container/queue)

---

### 📦 安装

```bash
go get github.com/yuancore/go-zen/container/queue
```

---

### 🚀 快速开始

#### 1. 客户端配置（ClientConfig）

| 参数            | 类型            | 默认值  | 描述                                                                 |
|-----------------|-----------------|---------|--------------------------------------------------------------------|
| `Addr`          | `string`        | 必填    | Redis 服务器地址（格式：`IP:Port`，如 `127.0.0.1:6379`）             |
| `Password`      | `string`        | `""`    | Redis 认证密码（空表示无密码）                                       |
| `DB`            | `int`           | `0`     | Redis 数据库编号（0-15）                                             |
| `PoolSize`      | `int`           | `20`    | 连接池大小（建议为最大预期并发数的 2 倍）                             |
| `DialTimeout`   | `time.Duration` | `10s`   | 建立连接的超时时间（如 `10 * time.Second`）                          |
| `ReadTimeout`   | `time.Duration` | `30s`   | 读取操作超时时间                                                     |
| `WriteTimeout`  | `time.Duration` | `30s`   | 写入操作超时时间                                                     |

**示例：初始化客户端**
```go
cfg := queue.ClientConfig{
    Addr:        "127.0.0.1:6379",
    Password:    "your_password",
    DB:          1,
    PoolSize:    50,
    DialTimeout: 15 * time.Second,
}
client := queue.NewClient(cfg, queue.WithLogger(zap.NewExample()))
defer client.Close()
```

#### 2. 服务端配置（ServiceConfig）

| 参数              | 类型               | 默认值           | 描述                                                                 |
|-------------------|--------------------|------------------|--------------------------------------------------------------------|
| `RedisAddress`    | `string`           | 必填             | Redis 地址（与客户端一致）                                           |
| `RedisPassword`   | `string`           | `""`             | Redis 密码                                                          |
| `RedisDB`         | `int`              | `1`              | Redis 数据库编号（默认与客户端区分）                                 |
| `Concurrency`     | `int`              | `10`             | 并发工作协程数（同时处理任务的最大数量）                              |
| `Queues`          | `map[string]int`   | `{"default": 1}` | 队列优先级配置（权重值越高优先级越高，如 `{"critical":5, "low":1}`） |
| `RetryStrategy`   | `RetryStrategy`    | `DefaultRetry`   | 自定义重试策略（需实现 `GetDelay` 方法）                             |
| `Logger`          | `*zap.Logger`      | `zap.NewNop()`   | 日志记录器（默认无日志输出）                                         |

**示例：服务端配置**
```go
cfg := queue.ServiceConfig{
    RedisAddress:  "127.0.0.1:6379",
    Concurrency:   30,
    Queues:        map[string]int{"high": 5, "default": 3, "low": 1},
    RetryStrategy: &CustomRetry{},
    Logger:        zap.NewExample(),
}
service := queue.NewService(&cfg)
```

#### 3. 任务选项（TaskOption）

| 选项函数                | 参数类型           | 默认值       | 描述                                                                 |
|-------------------------|--------------------|--------------|--------------------------------------------------------------------|
| `WithDelay(delay)`      | `time.Duration`    | `0`          | 延迟执行时间（如 `10 * time.Second`）                               |
| `WithMaxRetry(max)`     | `int`              | `0`          | 最大重试次数（`0` 表示不重试）                                      |
| `WithQueue(name)`       | `string`           | `"default"`  | 指定队列名称（需与服务端配置匹配）                                   |
| `WithTimeout(timeout)`  | `time.Duration`    | `0`          | 任务处理超时时间（超时后标记为失败）                                 |
| `WithDeadline(deadline)`| `time.Time`        | `time.Time{}`| 任务截止时间（超过时间不再执行）                                     |
| `WithUnique(ttl)`       | `time.Duration`    | `0`          | 唯一任务锁定时长（防止重复任务，如 `30 * time.Second`）              |

**示例：添加复杂任务**
```go
// 添加一个延迟5秒、最多重试3次、30秒内唯一的任务
info, err := client.Enqueue("task:process", payload,
    queue.WithDelay(5 * time.Second),
    queue.WithMaxRetry(3),
    queue.WithUnique(30 * time.Second),
    queue.WithQueue("high"),
)
```

---

### 🔧 高级用法

#### 1. 自定义重试策略
```go
type CustomRetry struct{}

func (r *CustomRetry) GetDelay(retryCount int, _ error, _ *asynq.Task) time.Duration {
    return time.Duration(retryCount) * 2 * time.Minute // 每次重试间隔加倍
}

// 注入服务端配置
cfg := queue.ServiceConfig{
    RetryStrategy: &CustomRetry{},
}
```

#### 2. 监控队列状态
```go
info, err := client.GetQueueInfo("high")
if err == nil {
    fmt.Printf("队列任务积压数: %d\n活跃Worker数: %d\n", info.Size, info.Active)
}
```

#### 3. 健康检查
```go
// 定时检查Redis连接
go func() {
    for {
        if err := client.HealthCheck(); err != nil {
            log.Printf("Redis连接异常: %v", err)
        }
        time.Sleep(30 * time.Second)
    }
}()
```

---

### ✨ 核心特性

| 特性                  | 说明                                                                 |
|-----------------------|--------------------------------------------------------------------|
| **Redis 集群支持**     | 支持单节点和集群模式                                                 |
| **任务优先级控制**     | 多队列权重分配，灵活调度高优先级任务                                 |
| **自动重试机制**       | 默认指数退避策略，支持自定义                                         |
| **线程安全设计**       | 单例客户端 + 读写锁，服务端协程池隔离                                |
| **唯一性任务保障**     | 基于 Redis 分布式锁，防止重复任务提交                                |

---

### ⚠️ 注意事项

1. **队列权重**  
   服务端的 `Queues` 配置中，权重值决定任务消费优先级（例如 `{"critical":5}` 表示 `critical` 队列处理速度是默认的 5 倍）。

2. **连接池大小**  
   `PoolSize` 建议设置为服务端 `Concurrency` 的 2 倍，避免连接竞争。

3. **唯一性任务**  
   使用 `WithUnique` 时，需确保所有 Redis 节点时间同步，防止锁提前失效。

4. **超时处理**  
   任务处理超时（`WithTimeout`）后会自动取消，需在处理器中处理上下文取消逻辑：
   ```go
   func handler(ctx context.Context, task *asynq.Task) error {
       select {
       case <-ctx.Done():
           return fmt.Errorf("任务超时")
       default:
           // 正常处理逻辑
       }
   }
   ```

---

### 🤝 参与贡献
[贡献指南](https://github.com/yuancore/go-zen/blob/main/CONTRIBUTING.md) | [提交Issue](https://github.com/yuancore/go-zen/issues)

---

## English

### 📖 Introduction

`queue` is a Redis-backed asynchronous task queue library for Go, supporting delayed tasks, unique jobs, and priority queues. Built on `asynq` with thread-safe client/server operations for distributed task scheduling.

GitHub URL: [github.com/yuancore/go-zen/container/queue](https://github.com/yuancore/go-zen/container/queue)

---

### 📦 Installation

```bash
go get github.com/yuancore/go-zen/container/queue
```

---

### 🚀 Quick Start

#### 1. Client Configuration (ClientConfig)

| Parameter         | Type             | Default      | Description                                                         |
|-------------------|------------------|--------------|---------------------------------------------------------------------|
| `Addr`            | `string`         | Required     | Redis server address (format: `IP:Port`, e.g., `127.0.0.1:6379`)    |
| `Password`        | `string`         | `""`         | Redis authentication password (empty for no auth)                   |
| `DB`              | `int`            | `0`          | Redis database index (0-15)                                         |
| `PoolSize`        | `int`            | `20`         | Connection pool size (recommended: 2x max concurrency)              |
| `DialTimeout`     | `time.Duration`  | `10s`        | Connection timeout (e.g., `10 * time.Second`)                       |
| `ReadTimeout`     | `time.Duration`  | `30s`        | Read operation timeout                                              |
| `WriteTimeout`    | `time.Duration`  | `30s`        | Write operation timeout                                             |

```go
cfg := queue.ClientConfig{
    Addr:        "127.0.0.1:6379",
    Password:    "your_password",
    DB:          1,
    PoolSize:    50,
    DialTimeout: 15 * time.Second,
}
client := queue.NewClient(cfg, queue.WithLogger(zap.NewExample()))
defer client.Close()
```

#### 2. Server Configuration (ServiceConfig)

| Parameter           | Type               | Default           | Description                                                         |
|---------------------|--------------------|-------------------|---------------------------------------------------------------------|
| `RedisAddress`      | `string`           | Required          | Redis server address                                                |
| `RedisPassword`     | `string`           | `""`              | Redis password                                                      |
| `RedisDB`           | `int`              | `1`               | Redis database index                                                |
| `Concurrency`       | `int`              | `10`              | Max concurrent workers                                              |
| `Queues`            | `map[string]int`   | `{"default": 1}`  | Queue priorities (higher weight = higher priority)                  |
| `RetryStrategy`     | `RetryStrategy`    | `DefaultRetry`    | Custom retry strategy (implement `GetDelay`)                        |
| `Logger`            | `*zap.Logger`      | `zap.NewNop()`    | Logger (no output by default)                                       |

```go
cfg := queue.ServiceConfig{
    RedisAddress:  "127.0.0.1:6379",
    Concurrency:   30,
    Queues:        map[string]int{"high": 5, "default": 3, "low": 1},
    RetryStrategy: &CustomRetry{},
    Logger:        zap.NewExample(),
}
service := queue.NewService(&cfg)
```

#### 3. Task Options (TaskOption)

| Option Function          | Parameter Type      | Default       | Description                                                         |
|--------------------------|---------------------|---------------|---------------------------------------------------------------------|
| `WithDelay(delay)`       | `time.Duration`     | `0`           | Delay task execution (e.g., `10 * time.Second`)                     |
| `WithMaxRetry(max)`      | `int`               | `0`           | Max retry attempts (`0` means no retry)                            |
| `WithQueue(name)`        | `string`            | `"default"`   | Target queue name (must match server config)                        |
| `WithTimeout(timeout)`   | `time.Duration`     | `0`           | Task processing timeout                                             |
| `WithDeadline(deadline)` | `time.Time`         | `time.Time{}` | Task deadline (no execution after this time)                        |
| `WithUnique(ttl)`        | `time.Duration`     | `0`           | Unique task lock TTL (e.g., `30 * time.Second`)                     |

```go
// Add a delayed task with retries and uniqueness
info, err := client.Enqueue("task:process", payload,
    queue.WithDelay(5 * time.Second),
    queue.WithMaxRetry(3),
    queue.WithUnique(30 * time.Second),
    queue.WithQueue("high"),
)
```

---

### 🔧 Advanced Usage

#### 1. Custom Retry Strategy
```go
type CustomRetry struct{}

func (r *CustomRetry) GetDelay(retryCount int, _ error, _ *asynq.Task) time.Duration {
    return time.Duration(retryCount) * 2 * time.Minute // Exponential backoff
}

// Inject into server config
cfg := queue.ServiceConfig{
    RetryStrategy: &CustomRetry{},
}
```

#### 2. Monitor Queue Status
```go
info, err := client.GetQueueInfo("high")
if err == nil {
    fmt.Printf("Pending tasks: %d\nActive workers: %d\n", info.Size, info.Active)
}
```

#### 3. Health Checks
```go
// Periodically check Redis connection
go func() {
    for {
        if err := client.HealthCheck(); err != nil {
            log.Printf("Redis connection error: %v", err)
        }
        time.Sleep(30 * time.Second)
    }
}()
```

---

### ✨ Key Features

| Feature                  | Description                                                     |
|--------------------------|-----------------------------------------------------------------|
| **Redis Cluster Support**| Single-node and cluster mode                                    |
| **Priority Queues**      | Weight-based task prioritization                                |
| **Auto Retry**           | Default exponential backoff, customizable strategies            |
| **Thread-Safe**          | Singleton client with RWMutex, worker isolation                 |
| **Unique Tasks**         | Redis-based lock to prevent duplicates                          |

---

### ⚠️ Important Notes

1. **Queue Weights**  
   Higher weights in `Queues` config mean higher priority (e.g., `{"high":5}` processes tasks 5x faster).

2. **Connection Pool**  
   Set `PoolSize` to 2x `Concurrency` to avoid contention.

3. **Unique Tasks**  
   Ensure Redis server time synchronization when using `WithUnique`.

4. **Timeout Handling**  
   Handle context cancellation in task handlers:
   ```go
   func handler(ctx context.Context, task *asynq.Task) error {
       select {
       case <-ctx.Done():
           return fmt.Errorf("task timeout")
       default:
           // Process task
       }
   }
   ```

---

### 🤝 Contributing
[Contribution Guide](https://github.com/yuancore/go-zen/blob/main/CONTRIBUTING.md) | [Open an Issue](https://github.com/yuancore/go-zen/issues)

[⬆ Back to Top](#中文)