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

type Request struct {
	Method  string
	URL     *url.URL
	Retries uint
	Client  *http.Client
}

func (r *Request) With(h http.Header, v url.Values, payload interface{}, ctx context.Context) (*http.Response, error) {
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
	var reqFunc funct.FailableNullaryFunc = func() error {
		if ctx != nil {
			if err := ctx.Err(); err != nil {
				return err
			}
		}
		response, err = r.Client.Do(req)
		return err
	}
	err = reqFunc.Retry(&funct.RetryConf{
		Retries:              r.Retries,
		ConcurrentErrHandler: true,
	})
	return response, err
}
