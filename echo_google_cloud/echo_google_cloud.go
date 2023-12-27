package echo_google_cloud

import (
	"io"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/webdevelop-pro/go-common/configurator"
	"github.com/webdevelop-pro/go-common/verser"
	logger "github.com/webdevelop-pro/go-logger"
)

const (
	errorTypeKey   = "@type"
	errorTypeValue = "type.googleapis.com/google.devtools.clouderrorreporting.v1beta1.ReportedErrorEvent"
)

type EchoGoogleCloud struct {
	skip       bool
	reqContext echo.Context
}

func (h EchoGoogleCloud) echoContextToServiceContext() logger.ServiceContext {
	serviceCtx := logger.ServiceContext{
		HttpRequest: &logger.HttpRequestContext{
			Method:             h.reqContext.Request().Method,
			URL:                h.reqContext.Request().URL.String(),
			UserAgent:          h.reqContext.Request().UserAgent(),
			Referrer:           h.reqContext.Request().Referer(),
			RemoteIp:           h.reqContext.Request().RemoteAddr,
			ResponseStatusCode: h.reqContext.Response().Status,
		},
	}
	if service := verser.GetService(); service != "" {
		serviceCtx.Service = service
	}

	if version := verser.GetVersion(); version != "" {
		serviceCtx.Version = version
	}

	repository := verser.GetRepository()
	revisionID := verser.GetRevisionID()
	if repository != "" || revisionID != "" {
		serviceCtx.SourceReference = &logger.SourceReference{
			Repository: repository,
			RevisionID: revisionID,
		}
	}

	if h.reqContext != nil {
		if user := h.reqContext.Get("user"); user != nil {
			serviceCtx.User = user.(string)
		}
	}
	return serviceCtx
}

func (h EchoGoogleCloud) Run(e *zerolog.Event, level zerolog.Level, s string) {
	if h.skip {
		return
	}

	switch level {
	case zerolog.ErrorLevel, zerolog.FatalLevel, zerolog.PanicLevel:
		if h.reqContext != nil {
			e.Interface("serviceContext", h.echoContextToServiceContext())
		}
		e.Str(errorTypeKey, errorTypeValue)
	}
}

// NewLogger return logger instance
func NewEchoGCLogger(component string, logLevel string, output io.Writer, c echo.Context) logger.Logger {
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
		Hook(EchoGoogleCloud{reqContext: c, skip: !skipGoogleHook}).
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

	return logger.Logger{l.Logger()}
}

// DefaultStdoutLogger return default logger instance
func DefaultStdoutLogger(logLevel string, c echo.Context) logger.Logger {
	return NewEchoGCLogger("default", logLevel, os.Stdout, c)
}

// NewComponentLogger return default logger instance with custom component
func NewComponentLogger(component string, c echo.Context) logger.Logger {
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
