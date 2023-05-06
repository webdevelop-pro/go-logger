package main

import (
	"net/http"

	"github.com/pkg/errors"

	"github.com/labstack/echo/v4"
	"github.com/webdevelop-pro/go-common/verser"
	logger "github.com/webdevelop-pro/go-logger/echo_google_cloud"
)

var ErrNotFound = errors.New("row not found")

func FindUser() error {
	return errors.Errorf("%s: %d", ErrNotFound, 123)
}

var (
	service    string
	version    string
	repository string
	revisionID string
)

func main() {
	verser.SetServiVersRepoRevis(service, version, repository, revisionID)
	e := echo.New()
	defaultLogger := logger.NewComponentLogger("main", nil) // logger without context
	e.GET("/", func(c echo.Context) error {
		err := FindUser()
		c.Set("user", "uuid4-1234")                     // set up user id to have it in logs
		log := logger.NewComponentLogger("get-func", c) // logger with get request context
		log.Error().Stack().Err(err).Msg("error while getting element")
		return c.String(http.StatusOK, "Hello, World!")
	})
	defaultLogger.Fatal().Err(e.Start(":1323")).Msg("echo went down")
}
