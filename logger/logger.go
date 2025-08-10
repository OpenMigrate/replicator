package logger

import (
	"io"
	"log/slog"
	"os"
)

type Options struct {
	Verbose bool
	File    string
	JSON    bool
}

var (
	instance *slog.Logger
	closer   func() error = func() error { return nil }
)

// Init creates the logger and stores it globally for Get().
func Init(opts Options) (*slog.Logger, func() error, error) {
	var writer io.Writer = os.Stdout
	var newCloser func() error = func() error { return nil }

	// If file path is provided, use file; otherwise use stdout
	if opts.File != "" {
		f, err := os.OpenFile(opts.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
		if err != nil {
			return nil, nil, err
		}
		writer = f
		newCloser = f.Close
	}

	// Determine log level based on verbose flag
	level := slog.LevelInfo
	if opts.Verbose {
		level = slog.LevelDebug
	}

	// Create appropriate handler based on JSON flag
	var handler slog.Handler
	handlerOpts := &slog.HandlerOptions{Level: level}
	
	if opts.JSON {
		handler = slog.NewJSONHandler(writer, handlerOpts)
	} else {
		handler = slog.NewTextHandler(writer, handlerOpts)
	}

	instance = slog.New(handler)
	closer = newCloser

	return instance, closer, nil
}

// Get returns the initialized logger.
func Get() *slog.Logger {
	if instance == nil {
		panic("logger not initialized. call logger.Init in main first")
	}
	return instance
}

// Close closes the underlying file if one was opened.
func Close() error {
	return closer()
}

