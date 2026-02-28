# antgo/i18n - Internationalization (i18n) Library / 国际化(i18n)库

[中文](#中文) | [English](#english)

---

## 中文

### 📖 简介

`antgo/i18n` 是一款基于Go语言的高效国际化（i18n）库，旨在为应用程序提供多语言支持。支持从文件加载语言包，缓存翻译结果，自动处理多语言切换，且具有高性能和低内存消耗。  
适用于需要处理多语言支持、日期时间格式化、复数规则等场景。

GitHub地址: [github.com/yuancore/go-zen/i18n](https://github.com/yuancore/go-zen/i18n)

### 📦 安装

```bash
go get github.com/yuancore/go-zen/i18n
```

### 🚀 快速开始

#### 初始化国际化模块
```go
package main

import (
	"fmt"
	"github.com/yuancore/go-zen/i18n"
)

func main() {
	// 配置国际化设置
	config := i18n.Config{
		DefaultLang:    "en",      // 默认语言 | Default language
		FallbackLang:   "zh-CN",   // 备用语言 | Fallback language
		TranslationsDir: "./translations", // 翻译文件目录 | Translation files directory
		SupportedLangs: []string{"en", "zh-CN", "es", "fr"}, // 支持的语言 | Supported languages
		CacheEnabled:   true,      // 启用翻译缓存 | Enable translation cache
		MaxCacheSize:   100,      // 缓存最大条目数 | Maximum cache entries
	}

	// 初始化国际化模块
	if err := i18n.New(config); err != nil {
		fmt.Println("初始化错误:", err)
		return
	}

	// 使用翻译功能
	fmt.Println(i18n.T(nil, "hello_world")) // 输出: Hello, World!
}
```

#### 使用翻译功能
```go
package main

import (
	"fmt"
	"github.com/yuancore/go-zen/i18n"
)

func main() {
	// 假设已初始化国际化模块

	// 获取翻译文本
	fmt.Println(i18n.T(nil, "hello_world")) // 输出: Hello, World!
	
	// 获取复数形式翻译
	fmt.Println(i18n.TPlural(nil, 2, "item_count", 2)) // 输出: 2 items

	// 获取本地化日期时间格式
	fmt.Println(i18n.TDate(nil, time.Now())) // 输出: 2025-02-08T12:00:00Z（根据语言设置可能不同）
}
```

### ✨ 核心特性

| 特性               | 描述                                                                 |
|--------------------|--------------------------------------------------------------------|
| **支持多语言**      | 通过加载不同语言的翻译包支持多语言 | Multi-language support via loading different language bundles |
| **高效缓存**        | 启用翻译缓存，减少重复翻译请求 | Translation caching for reducing repeated translation requests |
| **日期时间本地化**  | 根据语言设置提供本地化日期和时间格式 | Localized date and time formatting based on language settings |
| **复数规则支持**    | 自动处理单数和复数翻译 | Automatic handling of singular and plural translations |
| **严格的错误处理**  | 提供详细的错误信息，避免崩溃 | Safe error handling with detailed error reporting |

### ⚠️ 注意事项
1. 配置中的翻译文件目录必须包含有效的翻译文件（支持JSON, TOML, YAML格式）。
2. 确保输入的语言代码在支持的语言列表中有效。
3. 默认语言和备用语言必须存在于加载的语言包中。

### 🤝 参与贡献
[贡献指南](https://github.com/yuancore/go-zen/blob/main/CONTRIBUTING.md) | [提交Issue](https://github.com/yuancore/go-zen/issues)

---

## English

### 📖 Introduction

`antgo/i18n` is an efficient internationalization (i18n) library for Go, designed to provide multi-language support for applications. It supports loading language bundles from files, caching translation results, and automatically handling language switching with high performance and low memory consumption.  
Ideal for scenarios that require multi-language support, date/time formatting, pluralization rules, etc.

GitHub URL: [github.com/yuancore/go-zen/i18n](https://github.com/yuancore/go-zen/i18n)

### 📦 Installation

```bash
go get github.com/yuancore/go-zen/i18n
```

### 🚀 Quick Start

#### Initializing the Internationalization Module
```go
package main

import (
	"fmt"
	"github.com/yuancore/go-zen/i18n"
)

func main() {
	// Configuration for internationalization
	config := i18n.Config{
		DefaultLang:    "en",      // Default language
		FallbackLang:   "zh-CN",   // Fallback language
		TranslationsDir: "./translations", // Translation files directory
		SupportedLangs: []string{"en", "zh-CN", "es", "fr"}, // Supported languages
		CacheEnabled:   true,      // Enable translation cache
		MaxCacheSize:   100,      // Maximum cache size
	}

	// Initialize the internationalization module
	if err := i18n.New(config); err != nil {
		fmt.Println("Initialization error:", err)
		return
	}

	// Using translation feature
	fmt.Println(i18n.T(nil, "hello_world")) // Output: Hello, World!
}
```

#### Using the Translation Features
```go
package main

import (
	"fmt"
	"github.com/yuancore/go-zen/i18n"
)

func main() {
	// Assume the internationalization module has been initialized

	// Get translated text
	fmt.Println(i18n.T(nil, "hello_world")) // Output: Hello, World!
	
	// Get pluralized translation
	fmt.Println(i18n.TPlural(nil, 2, "item_count", 2)) // Output: 2 items

	// Get localized date and time format
	fmt.Println(i18n.TDate(nil, time.Now())) // Output: 2025-02-08T12:00:00Z (depending on language setting)
}
```

### ✨ Key Features

| Feature               | Description                                                             |
|-----------------------|-------------------------------------------------------------------------|
| **Multi-language Support** | Supports multiple languages by loading different language bundles |
| **Efficient Caching**  | Enables translation caching to reduce repeated translation requests |
| **Localized Date/Time** | Provides localized date and time formats based on language settings |
| **Pluralization Support** | Automatically handles singular and plural translations |
| **Safe Error Handling** | Provides detailed error messages to avoid crashes |

### ⚠️ Important Notes
1. The translation files directory specified in the configuration must contain valid translation files (supports JSON, TOML, YAML formats).
2. Ensure that the input language code is valid in the list of supported languages.
3. The default and fallback languages must exist within the loaded language bundles.

### 🤝 Contributing
[Contribution Guide](https://github.com/yuancore/go-zen/blob/main/CONTRIBUTING.md) | [Open an Issue](https://github.com/yuancore/go-zen/issues)

[⬆ Back to Top](#中文)
