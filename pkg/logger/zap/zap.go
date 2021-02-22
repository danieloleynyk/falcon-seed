package zap

import (
	"falcon-seed/pkg/config"
	"falcon-seed/pkg/logger"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"net/url"
	"time"
)

// Log represents SugaredLogger, If you want to increase the performance consider
// using Logger instead
type Logger struct {
	logger *zap.SugaredLogger
}

func (zap *Logger) Info(msg string, keysAndValues ...interface{}) {
	zap.logger.Infow(msg, keysAndValues...)
}

// New instantiates new SugaredLogger logger
func New(config *config.Logging) (*Logger, error) {
	var err error

	if err = zap.RegisterSink(LumberjackSinkName, func(u *url.URL) (zap.Sink, error) {
		return lumberjackSink{
			Logger: &lumberjack.Logger{
				Filename:   u.Opaque,
				MaxSize:    config.MaxSize,
				MaxBackups: config.MaxBackups,
				MaxAge:     config.MaxAge,
			},
		}, nil
	}); err != nil {
		return nil, err
	}

	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = SyslogTimeEncoder
	cfg.OutputPaths = append(cfg.OutputPaths, LumberjackOpaquePrefix+config.Path)

	zapLogger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	zap.ReplaceGlobals(zapLogger)

	defer zapLogger.Sync() // flushes buffer, if any
	return &Logger{
		logger: zapLogger.Sugar(),
	}, nil
}

func (zap *Logger) Debug(msg string, keysAndValues ...interface{}) {
	zap.logger.Debugw(msg, keysAndValues)
}

func (zap *Logger) Error(err error) {
	zap.logger.Error(err)
}

func (zap *Logger) Fatal(err error) {
	zap.logger.Fatal(err)
}

// Log logs using the initialized logger
func (zap *Logger) LogRequest(ctx echo.Context, source, msg string, err error, params map[string]interface{}) {

	if params == nil {
		params = make(map[string]interface{})
	}

	params[logger.Source] = source

	if id, ok := ctx.Get(logger.Id).(int); ok {
		params[logger.Id] = id
		params[logger.User] = ctx.Get(logger.User).(string)
	}

	if err != nil {
		params[logger.Error] = err
		zap.logger.Errorw(msg, params)
		return
	}

	zap.logger.Infof(msg, params)
}

func SyslogTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.UTC().String())
}

// GetLogger returns the global logger which was set to be the configured zap logger
func GetLogger() logger.Logger {
	return &Logger{
		logger: zap.S(),
	}
}
