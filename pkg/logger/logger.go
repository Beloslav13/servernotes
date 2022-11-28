package logger

import (
	"github.com/sirupsen/logrus"
)

var e *logrus.Entry

type Logger struct {
	*logrus.Entry
}

func GetLogger() Logger {
	return Logger{e}
}

func init() {
	logger := logrusConfig(logrus.New())
	e = logrus.NewEntry(logger)
}

func logrusConfig(l *logrus.Logger) *logrus.Logger {
	formatter := &logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	}

	l.SetFormatter(formatter)
	l.SetLevel(logrus.TraceLevel)
	return l
}
