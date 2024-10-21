package logger

type Option func(*Logger)

func WithEnabled(enabled bool) Option {
	return func(l *Logger) {
		l.enabled = enabled
	}
}
