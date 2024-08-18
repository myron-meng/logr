# logr

[TOC]

logr 是一个 Go `log/slog` 封装包，提供快速配置标准库 `log/slog` 以支持结构化日志打印，支持将日志输出到标准输出（stdout）和文件，并提供日志文件自动轮转功能。

## 特性

- 使用合理默认值，易于初始化，仅需一行代码即可完成日志配置
- 支持日志输出到标准输出（stdout）和文件
- 文件日志轮转功能
- 可自定义日志级别
- 可选择包含源代码位置信息
- 可自定义时间格式

## 安装

要使用此包，请运行：

```
go get github.com/myron-meng/logr
```

## 快速开始

### 基本用法（日志输出到标准输出）

如果您只需要将日志打印到标准输出，您可以在 `main()` 函数中使用一行代码初始化日志记录器：

```go
package main

import "github.com/myron-meng/logr"

func main() {
    logr.InitSlog() // 默认使用标准输出、日志级别为 slog.LevelInfo
    
    // 设置日志级别
    // logr.InitSlog(logr.WithLevel(slog.LevelInfo))

    // 根据函数设置日志级别
    // logr.InitSlog(logr.WithLevelFunc(func() slog.Level {
    //     l := slog.LevelDebug
    //     if config.Configs.Env == "prod" {
    //         l = slog.LevelInfo
    //     }
    //     return l
    // }))

    // 包含源代码位置信息
    // logr.InitSlog(logr.WithSource(true))

    // 自定义时间格式
    // logr.InitSlog(logr.WithTimeFormat(time.RFC3339))

    // 您的应用程序代码
}
```

### 日志输出到文件

要将日志输出到文件，请使用 `WithFileOutput` 选项和 `DefaultFileOutputConfig`：

```go
package main

import (
    "github.com/myron-meng/logr"
    "log/slog"
)

func main() {
    // 日志输出到相对路径，输出到当前目录下的 logs/app.log 文件
    logr.InitSlog(
        logr.WithFileOutput(logr.DefaultFileOutputConfig("logs/app.log")),
        logr.WithLevel(slog.LevelInfo),
    )

    // 或日志输出到绝对路径
    // logr.InitSlog(
    //     logr.WithFileOutput(logr.DefaultFileOutputConfig("/var/logs/myapp/app.log")),
    //     logr.WithLevel(slog.LevelInfo),
    // )
    
    // 您的应用程序代码
}
```

### 自定义文件输出

`FileOutputConfig` 结构允许您自定义文件日志记录的各个方面：

```go
logr.InitSlog(
    logr.WithFileOutput(logr.FileOutputConfig{
        Filename:   "logs/app.log",
        MaxSize:    100,  // 100 MB
        MaxBackups: 3,    // 保留 3 个旧文件
        MaxAge:     28,   // 28 天
        Compress:   true, // 压缩旧文件
    }),
    logr.WithLevel(slog.LevelDebug),
    logr.WithSource(true),
)
```

## 其他选项

- `WithLevel(level slog.Level)`: 设置固定的日志级别
- `WithLevelFunc(func() slog.Level)`: 设置动态日志级别
- `WithSource(bool)`: 在日志中包含源代码位置信息
- `WithTimeFormat(string)`: 自定义日志中的时间格式

## 使用日志记录器

初始化后，您可以使用标准的 `slog` 函数记录日志消息：

```go
slog.Info("应用程序已启动")
slog.Debug("这是一条调试消息")
slog.Error("发生错误", "error", err)

// 对于致命错误
logr.Fatal("发生致命错误", "error", err)
```

有关 `slog` 使用的更多详细信息，请参阅 Go 官方文档。