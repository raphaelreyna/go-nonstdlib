package apit

import (
	"github.com/raphaelreyna/go-nonstdlib/httpt"
	"net/http"
	"net/url"
	"path"
)

// Register adds a new service thats provided by one of the APIClients hosts.
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

// APIClient represents a client that has access to a slice of services
type APIClient struct {
	Client   *http.Client
	Services map[string]*httpt.Request
	hosts    map[string]*url.URL
}

// NewAPIClient creates returns a pointer to a new APIClient for the given hosts
// The hosts maps keys should be a user friendly name for the host.
func NewAPIClient(c *http.Client, hosts map[string]*url.URL) *APIClient {
	return &APIClient{
		Client:   c,
		hosts:    hosts,
		Services: map[string]*httpt.Request{},
	}
}
