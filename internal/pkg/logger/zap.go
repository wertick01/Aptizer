package logger

import (
	"encoding/json"
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapLogger struct {
	logger *zap.SugaredLogger
}

func NewZapLogger(cfg Config) (ILogger, error) {
	var zapCfg zap.Config

	zapRawCfg := fmt.Sprintf(`{
        "level": "%s",
        "encoding": "json",
        "outputPaths": ["stdout"],
        "errorOutputPaths": ["stderr"],
        "encoderConfig": {
            "messageKey": "message",
            "levelKey": "event_type",
            "levelEncoder": "capital",
            "timeKey": "datetime"
        }
    }`, cfg.Level)
	if cfg.CustomRawZapCfg != "" {
		zapRawCfg = cfg.CustomRawZapCfg
	}

	if err := json.Unmarshal([]byte(zapRawCfg), &zapCfg); err != nil {
		return nil, err
	}

	if zapCfg.EncoderConfig.EncodeTime == nil {
		zapCfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("02/Jan/2006:15:04:05 -0700")
	}

	l, err := zapCfg.Build()
	if err != nil {
		return nil, err
	}

	return &zapLogger{
		logger: l.Sugar().With(
			"version", cfg.Version,
			"host", cfg.Host,
			"context", cfg.Context,
		),
	}, nil
}

func (l *zapLogger) Debug(f *Fields) {
	l.with(f.Detail, f.UserID, f.Code).Debug(f.Message)
}

func (l *zapLogger) Info(f *Fields) {
	l.with(f.Detail, f.UserID, f.Code).Info(f.Message)
}

func (l *zapLogger) Warn(f *Fields) {
	l.with(f.Detail, f.UserID, f.Code).Warn(f.Message)
}

func (l *zapLogger) Error(f *Fields) {
	l.with(f.Detail, f.UserID, f.Code).Error(f.Message)
}

func (l *zapLogger) Fatal(f *Fields) {
	l.with(f.Detail, f.UserID, f.Code).Fatal(f.Message)
}

func (l *zapLogger) Panic(f *Fields) {
	l.with(f.Detail, f.UserID, f.Code).Panic(f.Message)
}

func (l *zapLogger) Sync() {
	l.logger.Sync()
}

func (l *zapLogger) with(d *Detail, userID, code string) *zap.SugaredLogger {
	return l.logger.With(
		"details", d,
		"_user_guid", userID,
		"code", code,
	)
}
