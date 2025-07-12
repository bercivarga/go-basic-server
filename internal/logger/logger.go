package logger

import (
	"log/slog"
	"os"
	"strings"
)

// New returns a configured slog.Logger and also sets it as the
// process-wide default (slog.SetDefault).
func New() *slog.Logger {
	opts := &slog.HandlerOptions{
		Level:       levelFromEnv(), // DEBUG | INFO | WARN | ERROR
		AddSource:   true,           // file:line
		ReplaceAttr: trimSrcPath,    // shorten long paths
	}

	var h slog.Handler
	if isTTY(os.Stdout.Fd()) {
		h = slog.NewTextHandler(os.Stdout, opts)
	} else {
		h = slog.NewJSONHandler(os.Stdout, opts)
	}

	l := slog.New(h)
	slog.SetDefault(l)
	return l
}

/* ---------- helpers ---------------------------------------------------- */

func levelFromEnv() slog.Level {
	switch strings.ToUpper(os.Getenv("LOG_LEVEL")) {
	case "DEBUG":
		return slog.LevelDebug
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func trimSrcPath(_ []string, a slog.Attr) slog.Attr {
	if a.Key != slog.SourceKey {
		return a
	}
	if src, ok := a.Value.Any().(*slog.Source); ok {
		if i := strings.LastIndex(src.File, "/"); i != -1 {
			src.File = src.File[i+1:]
		}
	}
	return a
}

func isTTY(fd uintptr) bool {
	fi, err := os.Stdout.Stat()
	return err == nil && (fi.Mode()&os.ModeCharDevice) != 0
}
