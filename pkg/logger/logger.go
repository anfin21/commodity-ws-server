package logger

import (
	"commodity-ws-server/configs"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type (
	Logger      = zap.SugaredLogger
	PlainLogger = zap.Logger
)

const (
	defaultLogLevel         = zapcore.InfoLevel
	minHighPriorityLogLevel = zapcore.ErrorLevel
)

func New(cfg *configs.Config) *Logger {
	var encoder zapcore.Encoder
	if cfg.Server.Env == configs.ServerEnvLocalhost {
		encoder = getLocalEncoder()
	} else {
		encoder = getGCPEncoder()
	}

	core := getCore(encoder)
	plain := zap.New(
		core,
		zap.AddCaller(),
		zap.ErrorOutput(zapcore.Lock(os.Stderr)),
	)

	sugared := plain.Sugar()
	return sugared
}

func getLocalEncoder() zapcore.Encoder {
	encoderConfig := localEncoderConfig()
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func localEncoderConfig() zapcore.EncoderConfig {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return encoderConfig
}

func getGCPEncoder() zapcore.Encoder {
	encoderConfig := gcpEncoderConfig()
	return newGCPEncoder(&encoderConfig)
}

func gcpEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "severity",
		NameKey:        "logger",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    encodeLevel(),
		EncodeTime:     zapcore.RFC3339TimeEncoder,
		EncodeDuration: zapcore.MillisDurationEncoder,
	}
}

func getCore(encoder zapcore.Encoder) zapcore.Core {
	return zapcore.NewTee(
		coreWritingLowPriorityLog(encoder),
		coreWritingHighPriorityLog(encoder),
	)
}

func coreWritingLowPriorityLog(encoder zapcore.Encoder) zapcore.Core {
	stdout := zapcore.Lock(os.Stdout) // lock for concurrent safe
	return zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(stdout),
		isLowPriorityLogLevel(),
	)
}

func isLowPriorityLogLevel() zapcore.LevelEnabler {
	// lowPriority used by info\debug\warn
	return zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= defaultLogLevel && lvl < minHighPriorityLogLevel
	})
}

func coreWritingHighPriorityLog(encoder zapcore.Encoder) zapcore.Core {
	stderr := zapcore.Lock(os.Stderr)
	return zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(stderr),
		isHighPriorityLogLevel(),
	)
}

func isHighPriorityLogLevel() zapcore.LevelEnabler {
	// highPriority used by error\panic\fatal
	return zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= defaultLogLevel && lvl >= minHighPriorityLogLevel
	})
}

func encodeLevel() zapcore.LevelEncoder {
	return func(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		switch l {
		case zapcore.DebugLevel:
			enc.AppendString("DEBUG")
		case zapcore.InfoLevel:
			enc.AppendString("INFO")
		case zapcore.WarnLevel:
			enc.AppendString("WARNING")
		case zapcore.ErrorLevel:
			enc.AppendString("ERROR")
		case zapcore.DPanicLevel:
			enc.AppendString("CRITICAL")
		case zapcore.PanicLevel:
			enc.AppendString("ALERT")
		case zapcore.FatalLevel:
			enc.AppendString("EMERGENCY")
		}
	}
}
