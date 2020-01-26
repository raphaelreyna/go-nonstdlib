package apit

import (
	"github.com/raphaelreyna/go-nonstdlib/httpt"
	"net/http"
	"net/url"
	"path"
)

func (s *APIClient) Register(name, host, endpoint, method string, retries uint) {
	u, _ := url.Parse(s.hosts[host].String())
	u.Path = path.Join(u.Path, endpoint)
	s.Services[name] = &httpt.Request{
		Method:  method,
		URL:     u,
		Retries: retries,
		Client:  s.Client,
	}
}

type APIClient struct {
	Client   *http.Client
	Services map[string]*httpt.Request
	hosts    map[string]*url.URL
}

func NewAPIClient(c *http.Client, hosts map[string]*url.URL) *APIClient {
	return &APIClient{
		Client: c,
		hosts:  hosts,
	}
}
