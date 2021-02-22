package zap

import "gopkg.in/natefinch/lumberjack.v2"

const (
	LumberjackOpaquePrefix = "lumberjack:"
	LumberjackSinkName     = "lumberjack"
)

// lumberjackSink is used as handler for zap rotating file logs
type lumberjackSink struct {
	*lumberjack.Logger
}

func (sink lumberjackSink) Sync() error {
	return nil
}
