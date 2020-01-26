package httpt

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestRequest(t *testing.T) {
	// Set up test server
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// Record method
			w.Header().Add("Method", r.Method)
			// Echo url values
			w.Header().Add("Test-Value", r.URL.Query().Get("Test-Value"))
			// Echo a test header
			w.Header().Add("Test-Header", r.Header.Get("Test-Header"))
			//Echo the body
			defer r.Body.Close()
			io.Copy(w, r.Body)
		}),
	)
	u, _ := url.Parse(server.URL)
	// Set up request to test server
	request := &Request{
		Method:  "GET",
		URL:     u,
		Retries: 1,
		Client:  server.Client(),
	}
	// Set up header, value, and body
	header := http.Header{}
	header.Add("Test-Header", "PASS")
	values := url.Values{}
	values.Add("Test-Value", "PASS")
	body := map[string]string{"test_body": "PASS"}

	var verify = func(response *http.Response) {
		// Check method
		method := response.Header.Get("Method")
		if method != request.Method {
			t.Errorf("method is: %s\nexpected: %s\n", method, request.Method)
		}
		// Check header
		tHeader := response.Header.Get("Test-Header")
		cHeader := header.Get("Test-Header")
		if tHeader != cHeader {
			t.Errorf("Test-Header is: %s\nexpected: %s\n", tHeader, cHeader)
		}
		// Check value
		tValue := response.Header.Get("Test-Value")
		cValue := values.Get("Test-Value")
		if tValue != cValue {
			t.Errorf("value is: %s\nexpected: %s\n", tValue, cValue)
		}
		// Check body
		tBody := map[string]string{}
		defer response.Body.Close()
		err := json.NewDecoder(response.Body).Decode(&tBody)
		if err != nil {
			t.Fatal(err)
		}
		if tBody["test_body"] != body["test_body"] {
			t.Errorf("body is: %s\nexpected: %s\n", tBody["test_body"], body["test_body"])
		}
	}

	var testWithoutContext = func() (*http.Response, error) {
		return request.With(nil, header, values, body)
	}

	response, err := testWithoutContext()
	if err != nil {
		t.Fatal(err)
	}
	verify(response)

	// Test with contexts if not a short test
	if !testing.Short() {
		// Background context
		var testWithBackgroundContext = func() (*http.Response, error) {
			ctx := context.Background()
			return request.With(ctx, header, values, body)
		}
		response, err := testWithBackgroundContext()
		if err != nil {
			t.Fatal(err)
		}
		verify(response)

		// Cancelled context
		var testWithCancelledContext = func() (*http.Response, error) {
			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			return request.With(ctx, header, values, body)
		}
		_, err = testWithCancelledContext()
		if err == nil {
			t.Fatalf("expected context cancelled error, got nil")
		}
	}
}
