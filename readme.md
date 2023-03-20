# Golang logging enhancements


This project adds enhancements to improve logging in go services. Service based on the [zerolog](https://github.com/rs/zerolog) and [pkg/errors](https://github.com/pkg/errors) with additional improvements:

- add stack trace for the error, panic and fatal errors
- add request context to every logging instance
- hook for seemless integration with [google cloud error reporting](https://cloud.google.com/error-reporting)

## Usage

```
package main

import (
  "net/http"

  "github.com/pkg/errors"
  "github.com/labstack/echo/v4"
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
    log := logger.NewComponentLogger("get-func", e) // logger with get request context
    log.Error().Stack().Err(err).Msg("log message with stack trace and context")
    return c.String(http.StatusOK, "Hello, World!")
  })
  e.Logger.Fatal(e.Start(":1323"))
}
```

Will output
```json
{
  "level": "error",
  "component": "get-func",
  "stack": [
    {
      "func": "errorFunc",
      "line": "12",
      "source": "main.go"
    },
    {
      "func": "main.func1",
      "line": "19",
      "source": "main.go"
    },
    {
      "func": "(*Echo).add.func1",
      "line": "575",
      "source": "echo.go"
    },
    {
      "func": "(*Echo).ServeHTTP",
      "line": "662",
      "source": "echo.go"
    },
    {
      "func": "serverHandler.ServeHTTP",
      "line": "2936",
      "source": "server.go"
    },
    {
      "func": "(*conn).serve",
      "line": "1995",
      "source": "server.go"
    },
    {
      "func": "goexit",
      "line": "1172",
      "source": "asm_arm64.s"
    }
  ],
  "error": "some error",
  "severity": "ERROR",
  "serviceContext": {
    "service": "",
    "version": "",
    "httpRequest": {
      "method": "GET",
      "url": "/",
      "userAgent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36",
      "referrer": "",
      "responseStatusCode": 0,
      "remoteIp": "127.0.0.1:62312"
    }
  },
  "@type": "type.googleapis.com/google.devtools.clouderrorreporting.v1beta1.ReportedErrorEvent",
  "time": "2023-03-20T18:08:42+01:00",
  "message": "log message with stack trace and context"
}
```

#### Notes


### Config

- `LOG_LEVEL` define log level, required
- `LOG_CONSOLE` add terminal colors, useful for local development. Default false
    
## Contributing
[TBD]
