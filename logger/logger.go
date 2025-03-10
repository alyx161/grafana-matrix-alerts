package logger

import (
	"fmt"
	"github.com/rs/zerolog"
	deflog "github.com/rs/zerolog/log"
	"os"
	"runtime"
	"strings"
	"time"
)

var Log zerolog.Logger

func InitLogger(debug bool) {
	writer := zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
		w.Out = os.Stdout
		w.TimeFormat = time.Stamp
	})
	Log = zerolog.New(writer).With().Timestamp().Logger()
	if !debug {
		Log = Log.Level(zerolog.InfoLevel)
	}

	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.CallerMarshalFunc = CallerWithFunctionName
	defaultCtxLog := Log.With().Bool("default_context_log", true).Caller().Logger()
	zerolog.DefaultContextLogger = &defaultCtxLog
	deflog.Logger = Log.With().Bool("global_log", true).Caller().Logger()
	// Swap trace and debug level colors so trace pops out the least
	zerolog.LevelColors[zerolog.TraceLevel] = 0
	zerolog.LevelColors[zerolog.DebugLevel] = 34 // blue
}

// CallerWithFunctionName is an implementation for zerolog.CallerMarshalFunc that includes the caller function name
// in addition to the file and line number.
func CallerWithFunctionName(pc uintptr, file string, line int) string {
	files := strings.Split(file, "/")
	file = files[len(files)-1]
	name := runtime.FuncForPC(pc).Name()
	fns := strings.Split(name, ".")
	name = fns[len(fns)-1]
	return fmt.Sprintf("%s:%d:%s()", file, line, name)
}
