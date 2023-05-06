package main

import (
	"github.com/pkg/errors"

	"github.com/webdevelop-pro/go-common/verser"
	"github.com/webdevelop-pro/go-logger"
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
	defaultLogger := logger.NewComponentLogger("main", nil) // logger without context
	err := FindUser()
	defaultLogger.Error().Stack().Err(err).Msg("error while getting element")
}
