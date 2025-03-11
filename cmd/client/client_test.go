package client

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeclineRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.EqualValues(t, req.URL.String(), "/api/v1/request/123/decline")
		assert.EqualValues(t, req.Header.Get("X-Api-Key"), "api-key-123")
		assert.EqualValues(t, req.Method, "POST")

		rw.Write([]byte(`OK`))
	}))

	defer server.Close()

	client := OverseerClientImpl{server.URL, "api-key-123", server.Client()}
	err := client.DeclineRequest("123")
	assert.Nil(t, err)
}

func TestDeclineRequest_WithError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.EqualValues(t, req.URL.String(), "/api/v1/request/123/decline")
		assert.EqualValues(t, req.Header.Get("X-Api-Key"), "api-key-123")
		assert.EqualValues(t, req.Method, "POST")

		http.Error(rw, "Internal Server Error", 500)
	}))

	defer server.Close()

	client := OverseerClientImpl{server.URL, "api-key-123", server.Client()}
	err := client.DeclineRequest("123")
	assert.Error(t, err)
}

func TestDeleteRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.EqualValues(t, req.URL.String(), "/api/v1/request/123")
		assert.EqualValues(t, req.Header.Get("X-Api-Key"), "api-key-123")
		assert.EqualValues(t, req.Method, "DELETE")

		rw.Write([]byte(`OK`))
	}))

	defer server.Close()

	client := OverseerClientImpl{server.URL, "api-key-123", server.Client()}
	err := client.DeleteRequest("123")
	assert.Nil(t, err)
}

func TestDeleteRequest_WithError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.EqualValues(t, req.URL.String(), "/api/v1/request/123")
		assert.EqualValues(t, req.Header.Get("X-Api-Key"), "api-key-123")
		assert.EqualValues(t, req.Method, "DELETE")

		http.Error(rw, "Internal Server Error", 500)
	}))

	defer server.Close()

	client := OverseerClientImpl{server.URL, "api-key-123", server.Client()}
	err := client.DeleteRequest("123")
	assert.Error(t, err)
}
