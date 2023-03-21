package logger

import "github.com/go-logr/logr"

var logger = logr.Discard()

func Get() logr.Logger {
	return logger
}

func Set(l logr.Logger) {
	logger = l
}
