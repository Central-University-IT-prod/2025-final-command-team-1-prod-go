package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *zap.Logger

func NewLogger() {
	// var level zapcore.Level
	// if err := level.UnmarshalText([]byte(cfg.Level)); err != nil {
	// 	return nil, fmt.Errorf("invalid log level: %w", err)
	// }

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	encoder := zapcore.NewJSONEncoder(encoderConfig)
	// switch cfg.Format {
	// case "json":
	// 	encoder = zapcore.NewJSONEncoder(encoderConfig)
	// case "text":
	// 	encoder = zapcore.NewConsoleEncoder(encoderConfig)
	// default:
	// 	return nil, fmt.Errorf("invalid log format: %w", cfg.Format)
	// }

	// cores := make([]zapcore.Core, 0, len(cfg.OutputsPaths))
	// for _, output := range cfg.OutputsPaths {
	// 	writer, err := getLogWriter(output)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	core := zapcore.NewCore(encoder, writer, level)
	// 	cores = append(cores, core)
	// }

	stdoutW, err := getLogWriter("stdout")
	if err != nil {
		panic("Failed create logger")
	}

	pathOut, err := getLogWriter("./log/app.log")
	if err != nil {
		panic("Failed create logger")
	}
	cores := []zapcore.Core{
		zapcore.NewCore(
			encoder,
			stdoutW,
			zapcore.ErrorLevel,
		),
		zapcore.NewCore(
			encoder,
			pathOut,
			zapcore.ErrorLevel,
		),
		zapcore.NewCore(
			encoder,
			stdoutW,
			zapcore.FatalLevel,
		),
		zapcore.NewCore(
			encoder,
			pathOut,
			zapcore.FatalLevel,
		),
		zapcore.NewCore(
			encoder,
			stdoutW,
			zapcore.InfoLevel,
		),
		zapcore.NewCore(
			encoder,
			pathOut,
			zapcore.InfoLevel,
		),
	}

	combinedCore := zapcore.NewTee(cores...)

	opts := []zap.Option{
		zap.AddCaller(),
	}
	opts = append(opts, zap.AddStacktrace(zap.ErrorLevel))

	Logger = zap.New(combinedCore, opts...)
}

func getLogWriter(output string) (zapcore.WriteSyncer, error) {
	switch output {
	case "stdout":
		return zapcore.AddSync(os.Stdout), nil
	case "stderr":
		return zapcore.AddSync(os.Stderr), nil
	default:
		return zapcore.AddSync(&lumberjack.Logger{
			Filename:   output,
			MaxSize:    100, //mb
			MaxBackups: 3,
			MaxAge:     30, //days
		}), nil
	}
}
