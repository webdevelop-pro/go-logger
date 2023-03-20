package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/webdevelop-pro/go-logger"
)

func errorFunc() error {
	return errors.New("some error")
}

func main() {
	e := echo.New()
	defaultLogger := logger.NewComponentLogger("main", nil) // logger without context
	e.GET("/", func(c echo.Context) error {
		err := errorFunc()
		log := logger.NewComponentLogger("get-func", c) // logger with get request context
		log.Error().Stack().Err(err).Msg("log message with stack trace and context")
		return c.String(http.StatusOK, "Hello, World!")
	})
	defaultLogger.Fatal().Err(e.Start(":1323")).Msg("echo went down")
}
