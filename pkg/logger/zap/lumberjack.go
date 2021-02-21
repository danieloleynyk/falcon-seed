package zap

import "gopkg.in/natefinch/lumberjack.v2"

const (
	LumberjackOpaquePrefix = "lumberjack:"
	LumberjackSinkName     = "lumberjack"
)

type lumberjackSink struct {
	*lumberjack.Logger
}

func (sink lumberjackSink) Sync() error {
	return nil
}
