package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"sync"
	"time"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorGray   = "\033[90m"
)

type ColorHandler struct {
	out   io.Writer
	mu    *sync.Mutex
	level slog.Level
}

func NewColorHandler(out io.Writer, level slog.Level) *ColorHandler {
	return &ColorHandler{
		out:   out,
		mu:    &sync.Mutex{},
		level: level,
	}
}

func (h *ColorHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level
}

func (h *ColorHandler) Handle(_ context.Context, r slog.Record) error {
	levelColor := colorCyan
	switch {
	case r.Level >= slog.LevelError:
		levelColor = colorRed
	case r.Level >= slog.LevelWarn:
		levelColor = colorYellow
	case r.Level >= slog.LevelInfo:
		levelColor = colorGreen
	}

	timeStr := r.Time.Format(time.Kitchen)

	msg := fmt.Sprintf("%s%s%s %s%-5s%s %s",
		colorGray, timeStr, colorReset,
		levelColor, r.Level.String(), colorReset,
		r.Message,
	)

	r.Attrs(func(a slog.Attr) bool {
		msg += fmt.Sprintf(" %s%s%s %v", colorCyan, a.Key, colorReset, a.Value)
		return true
	})

	msg += "\n"

	h.mu.Lock()
	defer h.mu.Unlock()
	_, err := fmt.Fprint(h.out, msg)
	return err
}

func (h *ColorHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *ColorHandler) WithGroup(name string) slog.Handler {
	return h
}
