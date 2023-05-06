package logger

import "github.com/rs/zerolog"

// Context is default context for the web request and response
type Context interface {
	// Request returns `*http.Request`.
	Request() interface{}

	// Response returns `*Response`.
	Response() interface{}

	RealIP() string

	// Get retrieves data from the context.
	Get(key string) interface{}
}

// Logger is wrapper struct around logger.Logger that adds some custom functionality
type Logger struct {
	zerolog.Logger
}

// ServiceContext contain info for all logs
type ServiceContext struct {
	Service         string              `json:"service"`
	Version         string              `json:"version"`
	User            string              `json:"user,omitempty"`
	HttpRequest     *HttpRequestContext `json:"httpRequest,omitempty"`
	SourceReference *SourceReference    `json:"sourceReference,omitempty"`
}

// SourceReference repositary name and revision id
type SourceReference struct {
	Repository string `json:"repository"`
	RevisionID string `json:"revisionId"`
}

// HttpRequestContext http request context
type HttpRequestContext struct {
	Method             string `json:"method"`
	URL                string `json:"url"`
	UserAgent          string `json:"userAgent"`
	Referrer           string `json:"referrer"`
	ResponseStatusCode int    `json:"responseStatusCode"`
	RemoteIp           string `json:"remoteIp"`
}
