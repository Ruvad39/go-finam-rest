package finam

import (
	"log/slog"
	"os"
)

// logger = slog
var (
	logLevel = &slog.LevelVar{} // INFO
	log      = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	})).With(slog.String("package", "go-alor-http"))
)

func SetLogger(logger *slog.Logger) {
	log = logger
}

// SetLogDebug установим уровень логирования Debug
func SetLogDebug(debug bool) {
	if debug {
		logLevel.Set(slog.LevelDebug)
	} else {
		logLevel.Set(slog.LevelInfo)
	}

}
