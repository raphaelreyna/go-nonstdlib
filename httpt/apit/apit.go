package apit

import (
	"github.com/raphaelreyna/go-nonstdlib/httpt"
	"net/http"
	"net/url"
	"path"
)

func (s *APIClient) Register(name, host, endpoint, method string, retries uint) error {
	u, err := url.Parse(s.hosts[host].String())
	if err != nil {
		return err
	}
	u.Path = path.Join(u.Path, endpoint)
	s.Services[name] = &httpt.Request{
		Method:  method,
		URL:     u,
		Retries: retries,
		Client:  s.Client,
	}
	return nil
}

type APIClient struct {
	Client   *http.Client
	Services map[string]*httpt.Request
	hosts    map[string]*url.URL
}

func NewAPIClient(c *http.Client, hosts map[string]*url.URL) {
	_ = &APIClient{
		Client:   c,
		hosts:    hosts,
		Services: map[string]*httpt.Request{},
	}
}
