package xcomp

import (
	"os"
	"runtime"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Debug(msg string, fields ...LogField)
	Info(msg string, fields ...LogField)
	Warn(msg string, fields ...LogField)
	Error(msg string, fields ...LogField)
	Fatal(msg string, fields ...LogField)
	Panic(msg string, fields ...LogField)

	With(fields ...LogField) Logger
	WithContext(key string, value any) Logger

	GetServiceName() string
}

type LogField struct {
	Key   string
	Value any
}

func Field(key string, value any) LogField {
	return LogField{Key: key, Value: value}
}

type ZapLogger struct {
	logger *zap.Logger
	sugar  *zap.SugaredLogger
}

func NewLogger(configService *ConfigService) Logger {
	return NewLoggerWithConfig(configService)
}

// isTerminal checks if the output is a terminal that supports colors
func isTerminal() bool {
	// Check if we're on Windows
	if runtime.GOOS == "windows" {
		// On Windows, check for TERM environment variable or if we're in ConEmu/Windows Terminal
		if os.Getenv("TERM") != "" || os.Getenv("WT_SESSION") != "" || os.Getenv("ConEmuPID") != "" {
			return true
		}
		return false
	}

	// On Unix-like systems, check if stdout is a terminal
	if fileInfo, err := os.Stdout.Stat(); err == nil {
		return (fileInfo.Mode() & os.ModeCharDevice) != 0
	}

	return false
}

// shouldUseColors determines if colors should be used based on terminal support and configuration
func shouldUseColors(configService *ConfigService, format string) bool {
	// Only use colors for console format
	if format != "console" && format != "text" {
		return false
	}

	// Check if colors are explicitly disabled
	if configService.GetBool("logging.disable_colors", false) {
		return false
	}

	// Check if colors are explicitly enabled
	if configService.GetBool("logging.force_colors", false) {
		return true
	}

	// Check environment variables
	if os.Getenv("NO_COLOR") != "" || os.Getenv("TERM") == "dumb" {
		return false
	}

	if os.Getenv("FORCE_COLOR") != "" || os.Getenv("CLICOLOR_FORCE") != "" {
		return true
	}

	// Auto-detect terminal support
	return isTerminal()
}

func NewLoggerWithConfig(configService *ConfigService) Logger {
	var config zap.Config

	// Determine if we should use development or production config
	isDevelopment := configService.GetBool("logging.development", false)
	if isDevelopment {
		config = zap.NewDevelopmentConfig()
	} else {
		config = zap.NewProductionConfig()
	}

	// Set log level
	level := configService.GetString("logging.level", "info")
	switch level {
	case "debug":
		config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "info":
		config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "warn":
		config.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "error":
		config.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	case "fatal":
		config.Level = zap.NewAtomicLevelAt(zapcore.FatalLevel)
	default:
		config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	}

	// Set output format
	format := configService.GetString("logging.format", "json")
	switch format {
	case "json":
		config.Encoding = "json"
	case "console", "text":
		config.Encoding = "console"
	default:
		config.Encoding = "json"
	}

	// Set output paths
	outputPaths := configService.GetString("logging.output_paths", "stdout")
	if outputPaths != "" {
		config.OutputPaths = []string{outputPaths}
	} else {
		config.OutputPaths = []string{"stdout"}
	}

	errorOutputPaths := configService.GetString("logging.error_output_paths", "stderr")
	if errorOutputPaths != "" {
		config.ErrorOutputPaths = []string{errorOutputPaths}
	} else {
		config.ErrorOutputPaths = []string{"stderr"}
	}

	// Configure encoder
	config.EncoderConfig.TimeKey = configService.GetString("logging.time_key", "timestamp")
	config.EncoderConfig.LevelKey = configService.GetString("logging.level_key", "level")
	config.EncoderConfig.MessageKey = configService.GetString("logging.message_key", "message")
	config.EncoderConfig.CallerKey = configService.GetString("logging.caller_key", "caller")
	config.EncoderConfig.StacktraceKey = configService.GetString("logging.stacktrace_key", "stacktrace")

	// Set time encoder
	timeFormat := configService.GetString("logging.time_format", "iso8601")
	switch timeFormat {
	case "iso8601":
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	case "rfc3339":
		config.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	case "epoch":
		config.EncoderConfig.EncodeTime = zapcore.EpochTimeEncoder
	case "millis":
		config.EncoderConfig.EncodeTime = zapcore.EpochMillisTimeEncoder
	default:
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}

	// Set level encoder with automatic color detection
	levelFormat := configService.GetString("logging.level_format", "capital")
	useColors := shouldUseColors(configService, format)

	switch levelFormat {
	case "capital":
		if useColors {
			config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		} else {
			config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		}
	case "lower":
		if useColors {
			config.EncoderConfig.EncodeLevel = zapcore.LowercaseColorLevelEncoder
		} else {
			config.EncoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
		}
	case "color":
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	default:
		if useColors {
			config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		} else {
			config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		}
	}

	// Set caller encoder with colors if supported
	if useColors {
		config.EncoderConfig.EncodeCaller = zapcore.FullCallerEncoder
	} else {
		config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	}

	// Enable caller if configured
	config.DisableCaller = !configService.GetBool("logging.enable_caller", true)
	config.DisableStacktrace = !configService.GetBool("logging.enable_stacktrace", false)

	logger, err := config.Build()
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}

	return &ZapLogger{
		logger: logger,
		sugar:  logger.Sugar(),
	}
}

func NewDevelopmentLogger() Logger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic("Failed to initialize development logger: " + err.Error())
	}

	return &ZapLogger{
		logger: logger,
		sugar:  logger.Sugar(),
	}
}

func (l *ZapLogger) GetServiceName() string {
	return "Logger"
}

func (l *ZapLogger) Debug(msg string, fields ...LogField) {
	l.logger.Debug(msg, l.convertFields(fields)...)
}

func (l *ZapLogger) Info(msg string, fields ...LogField) {
	l.logger.Info(msg, l.convertFields(fields)...)
}

func (l *ZapLogger) Warn(msg string, fields ...LogField) {
	l.logger.Warn(msg, l.convertFields(fields)...)
}

func (l *ZapLogger) Error(msg string, fields ...LogField) {
	l.logger.Error(msg, l.convertFields(fields)...)
}

func (l *ZapLogger) Fatal(msg string, fields ...LogField) {
	l.logger.Fatal(msg, l.convertFields(fields)...)
}

func (l *ZapLogger) Panic(msg string, fields ...LogField) {
	l.logger.Panic(msg, l.convertFields(fields)...)
}

func (l *ZapLogger) With(fields ...LogField) Logger {
	return &ZapLogger{
		logger: l.logger.With(l.convertFields(fields)...),
		sugar:  l.logger.Sugar(),
	}
}

func (l *ZapLogger) WithContext(key string, value any) Logger {
	return l.With(Field(key, value))
}

func (l *ZapLogger) convertFields(fields []LogField) []zap.Field {
	zapFields := make([]zap.Field, len(fields))
	for i, field := range fields {
		zapFields[i] = zap.Any(field.Key, field.Value)
	}
	return zapFields
}
