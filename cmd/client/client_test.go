package client

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func notOk(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func ok(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func equals(t *testing.T, a string, b string) {
	t.Helper()
	if a != b {
		t.Fatalf("'%s' does not equal '%s'", a, b)
	}
}

func TestDeclineRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		equals(t, req.URL.String(), "/api/v1/request/123/decline")
		equals(t, req.Header.Get("X-Api-Key"), "api-key-123")
		equals(t, req.Method, "POST")

		rw.Write([]byte(`OK`))
	}))

	defer server.Close()

	client := OverseerClient{server.URL, "api-key-123", server.Client()}
	err := client.DeclineRequest("123")
	ok(t, err)
}

func TestDeclineRequestError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		equals(t, req.URL.String(), "/api/v1/request/123/decline")
		equals(t, req.Header.Get("X-Api-Key"), "api-key-123")
		equals(t, req.Method, "POST")

		http.Error(rw, "Internal Server Error", 500)
	}))

	defer server.Close()

	client := OverseerClient{server.URL, "api-key-123", server.Client()}
	err := client.DeclineRequest("123")
	notOk(t, err)
}

func TestDeleteRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		equals(t, req.URL.String(), "/api/v1/request/123")
		equals(t, req.Header.Get("X-Api-Key"), "api-key-123")
		equals(t, req.Method, "DELETE")

		rw.Write([]byte(`OK`))
	}))

	defer server.Close()

	client := OverseerClient{server.URL, "api-key-123", server.Client()}
	err := client.DeleteRequest("123")
	ok(t, err)
}

func TestDeleteRequestError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		equals(t, req.URL.String(), "/api/v1/request/123")
		equals(t, req.Header.Get("X-Api-Key"), "api-key-123")
		equals(t, req.Method, "DELETE")

		http.Error(rw, "Internal Server Error", 500)
	}))

	defer server.Close()

	client := OverseerClient{server.URL, "api-key-123", server.Client()}
	err := client.DeleteRequest("123")
	notOk(t, err)
}
