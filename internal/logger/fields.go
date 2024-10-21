package logger

type Field struct {
	key   string
	value any
}

type fields []Field

func WithField(key string, value any) Field {
	return Field{key, value}
}

func WithError(err error) Field {
	return Field{"error", err.Error()}
}

func (f fields) args() []any {
	args := make([]any, 0, len(f)*2)
	for _, field := range f {
		args = append(args, field.key, field.value)
	}

	return args
}
