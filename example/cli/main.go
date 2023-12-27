package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"

	"github.com/webdevelop-pro/go-common/verser"
	logger "github.com/webdevelop-pro/go-logger"
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

type Context struct {
}

func (ctx Context) RealIP() string {
	return ""
}
func (ctx Context) Response() *echo.Response {
	return nil
}
func (ctx Context) Request() *http.Request {
	return nil
}
func (ctx Context) Get(key string) interface{} {
	if key == "user" {
		return "user-uuid"
	}
	return ""
}

func main() {
	verser.SetServiVersRepoRevis(service, version, repository, revisionID)
	ctx := Context{}
	defaultLogger := logger.NewComponentLogger("main-cli", ctx)
	err := FindUser()
	defaultLogger.Error().Stack().Err(err).Msg(ErrNotFound.Error())
}
