package log

import (
	"go.uber.org/zap"
)

var logger *zap.SugaredLogger

// init ...
func init() {
	l, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	logger = l.Sugar()
	logger.Infoln("Initialized logger")
}

// Log returns the logger.
func Log() *zap.SugaredLogger {
	return logger
}

// SetLevel sets the log level for the logger.
func SetLevel(level string) {
	l, err := zap.ParseAtomicLevel(level)
	if err != nil {
		panic(err)
	}
	logger = logger.WithOptions(zap.IncreaseLevel(l))
}
