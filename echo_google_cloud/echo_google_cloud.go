package echo_google_cloud

import (
	"context"
	"io"
	"os"

	"github.com/rs/zerolog"
	"github.com/webdevelop-pro/go-common/configurator"
	logger "github.com/webdevelop-pro/go-logger"
)

const (
	errorTypeKey   = "@type"
	errorTypeValue = "type.googleapis.com/google.devtools.clouderrorreporting.v1beta1.ReportedErrorEvent"
)

type EchoGoogleCloud struct {
	skip bool
}

func (h EchoGoogleCloud) Run(e *zerolog.Event, level zerolog.Level, s string) {
	if h.skip {
		return
	}

	switch level {
	case zerolog.ErrorLevel, zerolog.FatalLevel, zerolog.PanicLevel:
		e.Str(errorTypeKey, errorTypeValue)
	}
}

// NewLogger return logger instance
func NewEchoGCLogger(component string, logLevel string, output io.Writer, c context.Context) logger.Logger {
	level, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		level = zerolog.InfoLevel
	}

	// locally we don't need details for google cloud logging
	_, skipGoogleHook := output.(zerolog.ConsoleWriter)
	l := zerolog.
		New(output).
		Level(level).
		Hook(logger.SeverityHook{}).
		Hook(logger.ContextHook{}).
		Hook(EchoGoogleCloud{skip: !skipGoogleHook}).
		With().Timestamp()

	// if level == zerolog.DebugLevel || level == zerolog.TraceLevel {
	// l = l.Caller()
	// }

	if component != "" {
		l = l.Str("component", component)
	}

	if c != nil {
		l = l.Ctx(c)
	}

	if err != nil {
		ll := l.Logger()
		ll.Error().Err(err).Interface("level", logLevel).Msg("cannot parse log level, using default info")
	}

	return logger.Logger{l.Logger()}
}

// DefaultStdoutLogger return default logger instance
func DefaultStdoutLogger(logLevel string, c context.Context) logger.Logger {
	return NewEchoGCLogger("default", logLevel, os.Stdout, c)
}

// NewComponentLogger return default logger instance with custom component
func NewComponentLogger(component string, c context.Context) logger.Logger {
	conf := configurator.NewConfigurator()
	cfg := conf.New("logger", &logger.Config{}).(*logger.Config)

	var output io.Writer
	// Beautiful output
	if cfg.LogConsole {
		output = zerolog.NewConsoleWriter()
	} else {
		output = os.Stdout
	}

	return NewEchoGCLogger(component, cfg.LogLevel, output, c)
}
