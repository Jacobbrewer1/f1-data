package logging

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
)

func GeneralLogger(appName string) error {
	cfg := NewConfig(Name(appName))
	_, err := CommonLogger(cfg)
	if err != nil {
		return err
	}
	return nil
}

func GeneralTextLogger(appName string) error {
	cfg := NewConfig(Name(appName))
	// Check config.
	if cfg == nil {
		return errors.New("logging config is nil")
	}
	// Validate config.
	if cfg.appName == "" {
		return errors.New("app name is empty")
	}
	writer := os.Stdout
	_, err := CommonLoggerWithOptions(cfg, writer, slog.LevelDebug, false)
	if err != nil {
		return fmt.Errorf("error creating logger: %w", err)
	}
	return nil
}

// CommonLogger constructs a logging with default options.
func CommonLogger(cfg *Config) (*slog.Logger, error) {
	// Check config.
	if cfg == nil {
		return nil, errors.New("logging config is nil")
	}
	// Validate config.
	if cfg.appName == "" {
		return nil, errors.New("app name is empty")
	}
	writer := os.Stdout
	return CommonLoggerWithOptions(cfg, writer, slog.LevelDebug, true)
}

// CommonLoggerWithOptions constructs a logging with custom options.
func CommonLoggerWithOptions(cfg *Config, w io.Writer, minLevel slog.Level, logToJson bool) (*slog.Logger, error) {
	opts := slog.HandlerOptions{
		AddSource:   true,
		Level:       minLevel,
		ReplaceAttr: replaceAttrs,
	}

	logger := new(slog.Logger)
	if logToJson {
		logger = slog.New(slog.NewJSONHandler(w, &opts))
	} else {
		logger = slog.New(slog.NewTextHandler(w, &opts))
	}

	logger = logger.With(
		KeyAppName, cfg.appName,
	)

	slog.SetDefault(logger)

	return logger, nil
}

// replaceAttrs is a slog.HandlerOptions.ReplaceAttr function that replaces some attributes.
func replaceAttrs(_ []string, a slog.Attr) slog.Attr {
	switch a.Key {
	case slog.SourceKey:
		// Cut the source file to a relative path.
		v := strings.Split(a.Value.String(), "/")
		idx := len(v) - 2
		if idx < 0 {
			idx = 0
		}
		a.Value = slog.StringValue(strings.Join(v[idx:], "/"))

		// Remove any curly braces from the source file. This is needed for the logstash parser.
		a.Value = slog.StringValue(strings.ReplaceAll(a.Value.String(), "{", ""))
		a.Value = slog.StringValue(strings.ReplaceAll(a.Value.String(), "}", ""))
	}
	return a
}
