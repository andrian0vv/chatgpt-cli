package logger

import (
	"io"
	"log/slog"
)

type Logger struct {
	logger  *slog.Logger
	enabled bool
}

func New(w io.Writer, opts ...Option) *Logger {
	l := slog.New(
		NewColoredHandler(
			slog.NewTextHandler(
				w,
				&slog.HandlerOptions{
					Level: slog.LevelDebug,
				},
			),
		),
	)

	log := &Logger{
		logger: l,
	}

	for _, opt := range opts {
		opt(log)
	}

	return log
}

func (l *Logger) Debug(message string, f ...Field) {
	if l.enabled {
		l.logger.Debug(message, fields(f).args()...)
	}
}

func (l *Logger) Info(message string, f ...Field) {
	if l.enabled {
		l.logger.Info(message, fields(f).args()...)
	}
}

func (l *Logger) Error(message string, f ...Field) {
	if l.enabled {
		l.logger.Error(message, fields(f).args()...)
	}
}

func (l *Logger) Warn(message string, f ...Field) {
	if l.enabled {
		l.logger.Warn(message, fields(f).args()...)
	}
}
