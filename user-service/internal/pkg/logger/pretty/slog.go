package pretty

import (
	"context"
	"encoding/json"
	"github.com/fatih/color"
	"io"
	"log"
	"log/slog"
)

type HandlerOptions struct {
	SlogOpts *slog.HandlerOptions
}

type Handler struct {
	opts HandlerOptions
	slog.Handler
	l     *log.Logger
	attrs []slog.Attr
}

func (opts HandlerOptions) PrettyHandler(out io.Writer) *Handler {
	h := &Handler{
		Handler: slog.NewJSONHandler(out, opts.SlogOpts),
		l:       log.New(out, "", 0),
	}

	return h
}

func (h *Handler) Handle(_ context.Context, r slog.Record) error {
	level := r.Level.String() + ":"

	switch r.Level {
	case slog.LevelDebug:
		level = color.MagentaString(level)
	case slog.LevelInfo:
		level = color.BlueString(level)
	case slog.LevelWarn:
		level = color.YellowString(level)
	case slog.LevelError:
		level = color.RedString(level)
	}

	fields := make(map[string]interface{}, r.NumAttrs())

	r.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()

		return true
	})

	for _, a := range h.attrs {
		fields[a.Key] = a.Value.Any()
	}

	var bytes []byte
	var err error

	if len(fields) > 0 {
		bytes, err = json.MarshalIndent(fields, "", "  ")
		if err != nil {
			return err
		}
	}

	timeStr := r.Time.Format("[2006-01-02 15:04:05.000]")
	msg := color.CyanString(r.Message)

	h.l.Println(
		timeStr,
		level,
		msg,
		color.WhiteString(string(bytes)),
	)

	return nil
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &Handler{
		Handler: h.Handler,
		l:       h.l,
		attrs:   attrs,
	}
}

func (h *Handler) WithGroup(name string) slog.Handler {
	return &Handler{
		Handler: h.Handler.WithGroup(name),
		l:       h.l,
	}
}
