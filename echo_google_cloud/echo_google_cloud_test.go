package echo_google_cloud

import (
	"bytes"
	"context"
	"io"
	"os"
	"testing"

	"github.com/pkg/errors"

	"github.com/webdevelop-pro/go-common/tests"
	logger "github.com/webdevelop-pro/go-logger"
)

func ReadStdout(r *os.File, w *os.File) string {
	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()
	w.Close()
	out := <-outC

	return out
}

func testLog(t *testing.T, ctx context.Context, expected string, logF func(log logger.Logger)) {
	r, w, _ := os.Pipe()
	os.Stdout = w

	os.Setenv("LOG_LEVEL", "info")
	log := NewComponentLogger("test", nil)

	logF(log)

	actual := ReadStdout(r, w)

	tests.CompareJsonBody(t, []byte(actual), []byte(expected))
}

func TestLog_Info(t *testing.T) {
	ctx := context.WithValue(context.Background(), logger.ServiceContextInfo, logger.ServiceContext{
		Service:   "test-service",
		Version:   "v0.1",
		User:      "0001",
		RequestID: "asd-asd-asd",
		MSGID:     "ttt-tttt-ttt",
		HttpRequest: &logger.HttpRequestContext{
			Method:             "GET",
			URL:                "/test",
			UserAgent:          "testagent",
			Referrer:           "testReferrer",
			ResponseStatusCode: 200,
			RemoteIp:           "0.0.0.0",
		},
	})

	expected := `
		{
			"level": "info",
			"component": "test",
			"info": "details",
			"severity": "INFO",
			"serviceContext": {
				"service": "test-service",
				"version": "v0.1",
				"user": "0001",
				"request_id": "asd-asd-asd",
				"msg_id": "ttt-tttt-ttt",
				"httpRequest": {
					"method": "GET",
					"url": "/test",
					"userAgent": "testagent",
					"referrer": "testReferrer",
					"responseStatusCode": 200,
					"remoteIp": "0.0.0.0"
				}
			},
			"time": "%any%",
			"message": "test message"
		}
	`

	testLog(t, ctx, expected, func(log logger.Logger) {
		log.Info().Ctx(ctx).Str("info", "details").Msg("test message")
	})
}

func TestLog_ErrorWithoutStack(t *testing.T) {
	ctx := context.WithValue(context.Background(), logger.ServiceContextInfo, logger.ServiceContext{
		Service:   "test-service",
		Version:   "v0.1",
		User:      "0001",
		RequestID: "asd-asd-asd",
		MSGID:     "ttt-tttt-ttt",
		HttpRequest: &logger.HttpRequestContext{
			Method:             "GET",
			URL:                "/test",
			UserAgent:          "testagent",
			Referrer:           "testReferrer",
			ResponseStatusCode: 200,
			RemoteIp:           "0.0.0.0",
		},
	})

	expected := `
		{
			"level": "error",
			"component": "test",
			"error": "some critical error",
			"severity": "ERROR",
			"serviceContext": {
				"service": "test-service",
				"version": "v0.1",
				"user": "0001",
				"request_id": "asd-asd-asd",
				"msg_id": "ttt-tttt-ttt",
				"httpRequest": {
					"method": "GET",
					"url": "/test",
					"userAgent": "testagent",
					"referrer": "testReferrer",
					"responseStatusCode": 200,
					"remoteIp": "0.0.0.0"
				}
			},
			"time": "%any%",
			"message": "test message"
		}
	`

	testLog(t, ctx, expected, func(log logger.Logger) {
		log.Error().Ctx(ctx).Err(errors.New("some critical error")).Msg("test message")
	})
}

func TestLog_ErrorWithStack(t *testing.T) {
	ctx := context.WithValue(context.Background(), logger.ServiceContextInfo, logger.ServiceContext{
		Service:   "test-service",
		Version:   "v0.1",
		User:      "0001",
		RequestID: "asd-asd-asd",
		MSGID:     "ttt-tttt-ttt",
		HttpRequest: &logger.HttpRequestContext{
			Method:             "GET",
			URL:                "/test",
			UserAgent:          "testagent",
			Referrer:           "testReferrer",
			ResponseStatusCode: 200,
			RemoteIp:           "0.0.0.0",
		},
	})

	err := errors.Wrap(errors.New("error message"), "from error")

	expected := `
		{
			"level": "error",
			"component": "test",
			"stack": [
				{
					"func": "TestLog_ErrorWithStack",
					"line": "157",
					"source": "echo_google_cloud_test.go"
				},
				{
					"func": "tRunner",
					"line": "1689",
					"source": "testing.go"
				},
				{
					"func": "goexit",
					"line": "1695",
					"source": "asm_amd64.s"
				}
			],
			"error": "from error: error message",
			"severity": "ERROR",
			"serviceContext": {
				"service": "test-service",
				"version": "v0.1",
				"user": "0001",
				"request_id": "asd-asd-asd",
				"msg_id": "ttt-tttt-ttt",
				"httpRequest": {
					"method": "GET",
					"url": "/test",
					"userAgent": "testagent",
					"referrer": "testReferrer",
					"responseStatusCode": 200,
					"remoteIp": "0.0.0.0"
				}
			},
			"time": "%any%",
			"message": "test message"
		}
	`

	testLog(t, ctx, expected, func(log logger.Logger) {
		log.Error().Ctx(ctx).Stack().Err(err).Msg("test message")
	})
}
