package logger

import (
	"context"
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	l    *Logger
	once sync.Once
)

func NewLogger(logger *logrus.Logger) *Logger {
	once.Do(func() {
		l = &Logger{logger}
	})
	return l
}

func (l *Logger) ContextWithLogger(ctx context.Context) context.Context {
	return context.WithValue(ctx, Logger{}, l)
}

func FromContext(ctx context.Context) *Logger {
	if ll, ok := ctx.Value(Logger{}).(*Logger); ok {
		return ll
	}
	return &Logger{logrus.New()}
}

type Logger struct {
	*logrus.Logger
}
