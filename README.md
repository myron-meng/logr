# sloginit Package

sloginit is a Go `log/slog` wrapper that provides easy-to-use functions for initializing and configuring `log/slog` in your Go applications.

## Features

- Easy initialization with sensible defaults
- Support for logging to stdout and files
- File rotation capabilities
- Customizable log levels
- Option to include source code location in logs
- Customizable time format

## Installation

To use this package, run:

```
go get github.com/myron-meng/sloginit
```

## Quick Start

### Basic Usage (Logging to stdout)

If you only need to print logs to stdout, you can initialize the logger with a single line in your `main()` function:

```go
package main

import "github.com/myron-meng/sloginit"

func main() {
    sloginit.Init()

    // or set log level
    // sloginit.InitSlog(sloginit.WithLevel(slog.LevelInfo))

    // or set log level by function
    // sloginit.InitSlog(sloginit.WithLevelFunc(func() slog.Level {
    //     l := slog.LevelDebug
    //     if os.Getenv("ENV") == "prod" {
    //         l = slog.LevelInfo
    //     }
    //     return l
    // }))

    // Your application code here
}
```

### Logging to a File

To log to a file, use the `WithFileOutput` option with `DefaultFileOutputConfig`:

```go
package main

import (
    "github.com/myron-meng/sloginit"
    
    "log/slog"
)

func main() {
    // Log to a relative path
    sloginit.InitSlog(
        sloginit.WithFileOutput(sloginit.DefaultFileOutputConfig("logs/app.log")),
        sloginit.WithLevel(slog.LevelInfo),
    )
    defer sloginit.Close()

    // Or log to an absolute path
    // sloginit.InitSlog(
    //     sloginit.WithFileOutput(sloginit.DefaultFileOutputConfig("/var/logs/myapp/app.log")),
    //     sloginit.WithLevel(slog.LevelInfo),
    // )
    
    // Your application code here
}
```

### Customizing File Output

The `FileOutputConfig` struct allows you to customize various aspects of file logging:

```go
sloginit.InitSlog(
    sloginit.WithFileOutput(sloginit.FileOutputConfig{
        Filename:   "logs/app.log",
        MaxSize:    100,  // 100 MB
        MaxBackups: 3,    // Keep 3 old files
        MaxAge:     28,   // 28 days
        Compress:   true, // Compress old files
    }),
    sloginit.WithLevel(slog.LevelDebug),
    sloginit.WithSource(true),
)
defer sloginit.Close()
```

## Additional Options

- `WithLevel(level slog.Level)`: Set a fixed log level
- `WithLevelFunc(func() slog.Level)`: Set level based on a function
- `WithSource(bool)`: Include source code location in logs
- `WithTimeFormat(string)`: Customize the time format in logs

## Using the Logger

After initialization, you can use the standard `slog` functions to log messages:

```go
slog.Info("Application started")
slog.Debug("This is a debug message")
slog.Error("An error occurred", "error", err)

// For fatal errors
sloginit.Fatal("A fatal error occurred", "error", err)
// or
slog.Log(context.TODO(),  slog.Level(12), msg, args...)
```

For more detailed information on `slog` usage, refer to the official Go documentation.