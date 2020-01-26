package apit

import (
	"net/http"
	"net/url"
	"testing"
)

func TestClient(t *testing.T) {
	// Create host url
	u, _ := url.Parse("http://test:1234/foo")
	// Create new client
	c := NewAPIClient(
		http.DefaultClient,
		map[string]*url.URL{
			"test_host": u,
		})
	// Register new endpoint
	err := c.Register(
		"test_service", "test_host",
		"bar/baz", "GET", 1,
	)
	if err != nil {
		t.Fatal(err)
	}

	// Check new url
	tu := c.Services["test_service"].URL.String()
	cu := u.String() + "/bar/baz"
	if tu != cu {
		t.Errorf("service url is: %s\nexpected: %s", tu, cu)
	}
}
