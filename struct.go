package logger

import (
	"github.com/rs/zerolog"
)

const (
	ServiceContextInfo = "service_context_info"
)

// Logger is wrapper struct around logger.Logger that adds some custom functionality
type Logger struct {
	zerolog.Logger
}

// ServiceContext contain info for all logs
type ServiceContext struct {
	Service         string              `json:"service"`
	Version         string              `json:"version"`
	User            string              `json:"user,omitempty"`
	RequestID       string              `json:"request_id,omitempty"`
	MSGID           string              `json:"msg_id,omitempty"`
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
