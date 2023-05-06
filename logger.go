package logger

import (
	"io"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"

	"github.com/webdevelop-pro/go-common/configurator"
)

func init() {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
}

// Printf is implementation of fx.Printer
func (l Logger) Printf(s string, args ...interface{}) {
	l.Info().Msgf(s, args...)
}

// NewLogger return logger instance
func NewLogger(component string, logLevel string, output io.Writer, c Context) Logger {
	level, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		level = zerolog.InfoLevel
	}

	l := zerolog.
		New(output).
		Level(level).
		Hook(SeverityHook{}).
		With().Stack().Timestamp()

	// if level == zerolog.DebugLevel || level == zerolog.TraceLevel {
	// l = l.Caller()
	// }

	if component != "" {
		l = l.Str("component", component)
	}

	if err != nil {
		ll := l.Logger()
		ll.Error().Err(err).Interface("level", logLevel).Msg("cannot parse log level, using default info")
	}

	return Logger{l.Logger()}
}

// DefaultStdoutLogger return default logger instance
func DefaultStdoutLogger(logLevel string, c Context) Logger {
	return NewLogger("default", logLevel, os.Stdout, c)
}

// NewComponentLogger return default logger instance with custom component
func NewComponentLogger(component string, c Context) Logger {
	conf := configurator.NewConfigurator()
	cfg := conf.New("logger", &Config{}).(*Config)

	var output io.Writer
	// Beautiful output
	if cfg.LogConsole {
		output = zerolog.NewConsoleWriter()
	} else {
		output = os.Stdout
	}

	return NewLogger(component, cfg.LogLevel, output, c)
}
