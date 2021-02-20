package zap

import (
	"falcon-seed/pkg/config"
	"falcon-seed/pkg/logger"
	"fmt"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"time"
)

var lumberjackLogger *lumberjack.Logger

// Log represents SugaredLogger, If you want to increase the performance consider
// using Logger instead
type Logger struct {
	logger *zap.SugaredLogger
}

// New instantiates new SugaredLogger logger
func New(config *config.Logging) *Logger {
	lumberjackLogger = &lumberjack.Logger{
		Filename:   config.Path,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
	}

	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = SyslogTimeEncoder
	zapLogger, _ := cfg.Build(zap.Hooks(lumberjackZapHook))

	defer zapLogger.Sync() // flushes buffer, if any
	return &Logger{
		logger: zapLogger.Sugar(),
	}
}

func (zap *Logger) Info(msg string) {
	zap.logger.Info(msg)
}

func (zap *Logger) Debug(msg string) {
	zap.logger.Debug(msg)
}

func (zap *Logger) Error(err error) {
	zap.logger.Error(err)
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

func lumberjackZapHook(e zapcore.Entry) error {
	lumberjackLogger.Write([]byte(fmt.Sprintf("%+v\n", e)))
	return nil
}
