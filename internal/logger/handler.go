package logger

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/fatih/color"
)

type ColoredHandler struct {
	handler slog.Handler
}

func NewColoredHandler(handler slog.Handler) *ColoredHandler {
	return &ColoredHandler{handler: handler}
}

func (h *ColoredHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *ColoredHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h.handler.WithAttrs(attrs)
}

func (h *ColoredHandler) WithGroup(name string) slog.Handler {
	return h.handler.WithGroup(name)
}

func (h *ColoredHandler) Handle(_ context.Context, record slog.Record) error {
	var outputColor *color.Color
	switch record.Level {
	case slog.LevelDebug:
		outputColor = color.New(color.FgHiBlack)
	case slog.LevelInfo:
		outputColor = color.New(color.FgHiCyan)
	case slog.LevelWarn:
		outputColor = color.New(color.FgHiYellow)
	case slog.LevelError:
		outputColor = color.New(color.FgHiRed)
	default:
		outputColor = color.New(color.FgHiWhite)
	}

	logData := map[string]interface{}{
		"time":    time.Now().Format(time.RFC3339),
		"level":   record.Level.String(),
		"message": record.Message,
	}

	record.Attrs(func(attr slog.Attr) bool {
		logData[attr.Key] = attr.Value.Any()
		return true
	})

	jsonLog, err := json.Marshal(logData)
	if err != nil {
		return err
	}

	_, err = outputColor.Println(string(jsonLog))

	return err
}
