package httprequest

import (
	"net/http"
	"time"
)

// IhttpRequest Interface implementation for custom http request.This helps in mock testing.
type IhttpRequest interface {
	RetryDo(req *http.Request, maxRetries int, timeout time.Duration,
		backoffStrategy []time.Duration) (*http.Response, error)
	PerformPost(args RequestDataParams, backOffStrategy []time.Duration, headers map[string]string) (*http.Response, error)
}

type HttpRequest struct {
}

func New() IhttpRequest {
	return HttpRequest{}
}

type RequestDataParams struct {
	Endpoint string            `json:"endpoint"`
	Data     []byte            `json:"data"`
	Params   map[string]string `json:"params"`
}
