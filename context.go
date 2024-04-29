package logger

import (
	"github.com/rs/zerolog"
	"github.com/webdevelop-pro/go-common/context/keys"
)

type ContextHook struct{}

func (h ContextHook) Run(e *zerolog.Event, level zerolog.Level, _ string) {
	ctx := e.GetCtx()

	serviceCtx, _ := keys.GetCtxValue(ctx, keys.LogInfo).(ServiceContext)

	e.Interface("serviceContext", serviceCtx)
}
