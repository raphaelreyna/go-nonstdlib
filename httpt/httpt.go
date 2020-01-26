package httpt

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/raphaelreyna/go-nonstdlib/funct"
	"io"
	"net/http"
	"net/url"
)

// Request holds parameters that are common between subsequent http requests to the same endpoint.
type Request struct {
	Method  string
	URL     *url.URL
	Retries uint
	Client  *http.Client
}

// With has the requests client make the request with the given parameters.
func (r *Request) With(ctx context.Context, h http.Header, v url.Values, payload interface{}) (*http.Response, error) {
	var (
		body     io.Reader
		response *http.Response
		err      error
	)
	if payload != nil {
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(payloadBytes)
	}
	url := r.URL.String()
	if v != nil {
		url = url + v.Encode()
	}
	req, err := http.NewRequest(r.Method, url, body)
	if err != nil {
		return nil, err
	}
	if h != nil {
		req.Header = h
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}
	var reqFunc = func() error {
		response, err = r.Client.Do(req)
		return err
	}
	conf := &funct.RetryConf{
		Retries:              r.Retries,
		ConcurrentErrHandler: true,
	}
	err = funct.Retry(conf, reqFunc)
	return response, err
}
