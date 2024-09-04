package sloginit

import (
	"context"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	fileLogger *lumberjack.Logger
	closeOnce  sync.Once
)

// Option is a function type used to modify slog.HandlerOptions
type Option func(*slog.HandlerOptions)

// LevelFunc is a function type that returns slog.Level
type LevelFunc func() slog.Level

// FileOutputConfig defines the configuration for file output
type FileOutputConfig struct {
	// Supports both relative and absolute paths. A filename starting with "/" is considered an absolute path,
	// otherwise it's relative. Relative paths are relative to the current working directory,
	// e.g., "logs/app.log" will be relative to the current working directory
	Filename string
	// Maximum size in megabytes
	MaxSize int
	// Maximum number of old log files to retain
	MaxBackups int
	// Maximum number of days to retain old log files
	MaxAge int
	// Whether to compress old files
	Compress bool
}

// defaultFileOutputConfig provides the default file output configuration
var defaultFileOutputConfig = FileOutputConfig{
	MaxSize:    100,  // 100 MB
	MaxBackups: 0,    // No limit
	MaxAge:     28,   // 28 days
	Compress:   true, // Compress by default
}

// DefaultFileOutputConfig returns the default file output configuration with the specified filename
func DefaultFileOutputConfig(filename string) FileOutputConfig {
	config := defaultFileOutputConfig
	config.Filename = filename
	return config
}

// FatalLevel defines the Fatal log level
const FatalLevel = slog.Level(12) // Set higher than Error (8)

// Init initializes slog
func Init(options ...Option) {
	opts := &slog.HandlerOptions{}
	var output io.Writer = os.Stdout // Default output to standard output

	// Apply all options
	for _, option := range options {
		option(opts)
	}

	// Create handler
	handler := slog.NewJSONHandler(output, opts)

	// Set global logger
	slog.SetDefault(slog.New(handler))
}

// Fatal logs a Fatal level message and terminates the program
func Fatal(msg string, args ...any) {
	slog.Log(context.TODO(), FatalLevel, msg, args...)
	os.Exit(1)
}

// WithLevel sets a fixed log level
func WithLevel(level slog.Level) Option {
	return func(opts *slog.HandlerOptions) {
		opts.Level = level
	}
}

// WithLevelFunc sets a dynamic log level
func WithLevelFunc(levelFunc LevelFunc) Option {
	return func(opts *slog.HandlerOptions) {
		opts.Level = levelFunc()
	}
}

// WithOutput sets the output destination
func WithOutput(w io.Writer) Option {
	return func(opts *slog.HandlerOptions) {
		handler := slog.NewJSONHandler(w, opts)
		slog.SetDefault(slog.New(handler))
	}
}

// WithFileOutput sets file output and configures log rotation
func WithFileOutput(config FileOutputConfig) Option {
	return func(opts *slog.HandlerOptions) {
		filename := config.Filename
		if !filepath.IsAbs(filename) {
			currentDir, err := os.Getwd()
			if err != nil {
				panic("Unable to get current working directory: " + err.Error())
			}
			filename = filepath.Join(currentDir, filename)
		}

		// Ensure log directory exists
		logDir := filepath.Dir(filename)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			panic("Unable to create log directory: " + err.Error())
		}

		fileLogger = &lumberjack.Logger{
			Filename:   filename,
			MaxSize:    config.MaxSize,
			MaxBackups: config.MaxBackups,
			MaxAge:     config.MaxAge,
			Compress:   config.Compress,
		}
		handler := slog.NewJSONHandler(fileLogger, opts)
		slog.SetDefault(slog.New(handler))
	}
}

// Close safely closes the log file if file output is being used
func Close() {
	closeOnce.Do(func() {
		if fileLogger != nil {
			_ = fileLogger.Close()
		}
	})
}

// WithSource sets whether to include source code location
func WithSource(addSource bool) Option {
	return func(opts *slog.HandlerOptions) {
		opts.AddSource = addSource
	}
}

// WithTimeFormat sets the time format
func WithTimeFormat(format string) Option {
	return func(opts *slog.HandlerOptions) {
		opts.ReplaceAttr = func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.String(a.Key, a.Value.Time().Format(format))
			}
			return a
		}
	}
}
