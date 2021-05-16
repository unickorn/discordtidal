package log

import "go.uber.org/zap"

var logger *zap.SugaredLogger

func Init() {
	l, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	logger = l.Sugar()
	logger.Info("initialized logger")
}

func Log() *zap.SugaredLogger {
	return logger
}
