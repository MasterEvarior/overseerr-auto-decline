package handler

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type OverseerClientMock struct {
	mock.Mock
}

func (m *OverseerClientMock) DeclineRequest(requestID string) error {
	m.Called()

	return nil
}

func (m *OverseerClientMock) DeleteRequest(requestID string) error {
	m.Called()

	return nil
}

func TestWebhookHandler_InvalidMethod(t *testing.T) {
	invalidMethods := []string{"GET", "PUT", "HEAD", "DELETE", "PATCH"}

	for _, method := range invalidMethods {
		handler := &Handler{}
		req := httptest.NewRequest(method, "/", nil)
		rec := httptest.NewRecorder()

		handler.WebhookHandler(rec, req)

		assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
	}
}

func TestWebhookHandler_BadRequestBody(t *testing.T) {
	invalidBodies := []string{`{}`, `{ "tmdbid": "321", "tvdbid": "123" }`}

	for _, body := range invalidBodies {
		handler := &Handler{}
		req := httptest.NewRequest("POST", "/", bytes.NewBuffer([]byte(body)))
		rec := httptest.NewRecorder()

		handler.WebhookHandler(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
}

func TestWebhookHandler_NoBannedMediaIDs(t *testing.T) {
	mockClient := &OverseerClientMock{}
	handler := &Handler{
		OverseerrClient: mockClient,
		BannedMediaIDs:  []string{"123"},
		DeleteRequests:  true,
	}
	payload := `{"request_id": "1", "tmdbid": "999", "tvdbid": "888"}`
	req := httptest.NewRequest("POST", "/", bytes.NewBuffer([]byte(payload)))
	rec := httptest.NewRecorder()

	handler.WebhookHandler(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockClient.AssertNotCalled(t, "DeclineRequest", "1")
	mockClient.AssertNotCalled(t, "DeleteRequest", "1")
	mockClient.AssertExpectations(t)
}

func TestWebhookHandler_DeclineRequest_Success(t *testing.T) {
	mockClient := &OverseerClientMock{}
	mockClient.On("DeclineRequest", mock.AnythingOfType("string")).Return(nil)
	handler := &Handler{
		OverseerrClient: mockClient,
		BannedMediaIDs:  []string{"999"},
		DeleteRequests:  false,
	}
	payload := `{"request_id": "1", "tmdbid": "999", "tvdbid": "888"}`
	req := httptest.NewRequest("POST", "/", bytes.NewBuffer([]byte(payload)))
	rec := httptest.NewRecorder()

	handler.WebhookHandler(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
	mockClient.AssertCalled(t, "DeclineRequest", "1")
	mockClient.AssertNotCalled(t, "DeleteRequest", "1")
	mockClient.AssertExpectations(t)
}

func TestWebhookHandler_DeclineRequest_Failure(t *testing.T) {
	mockClient := &OverseerClientMock{}
	mockClient.On("DeclineRequest", mock.AnythingOfType("string")).Return(errors.New("Error while declining request"))
	handler := &Handler{
		OverseerrClient: mockClient,
		BannedMediaIDs:  []string{"999"},
		DeleteRequests:  false,
	}
	payload := `{"request_id": "1", "tmdbid": "999", "tvdbid": "888"}`
	req := httptest.NewRequest("POST", "/", bytes.NewBuffer([]byte(payload)))
	rec := httptest.NewRecorder()

	handler.WebhookHandler(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	mockClient.AssertCalled(t, "DeclineRequest", "1")
	mockClient.AssertNotCalled(t, "DeleteRequest", "1")
	mockClient.AssertExpectations(t)
}

func TestWebhookHandler_DeleteRequest_Success(t *testing.T) {
	mockClient := &OverseerClientMock{}
	mockClient.On("DeclineRequest", mock.AnythingOfType("string")).Return(nil)
	mockClient.On("DeleteRequest", mock.AnythingOfType("string")).Return(nil)
	handler := &Handler{
		OverseerrClient: mockClient,
		BannedMediaIDs:  []string{"999"},
		DeleteRequests:  true,
	}
	payload := `{"request_id": "1", "tmdbid": "999", "tvdbid": "888"}`
	req := httptest.NewRequest("POST", "/", bytes.NewBuffer([]byte(payload)))
	rec := httptest.NewRecorder()

	handler.WebhookHandler(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
	mockClient.AssertCalled(t, "DeclineRequest", "1")
	mockClient.AssertCalled(t, "DeleteRequest", "1")
	mockClient.AssertExpectations(t)
}

func TestWebhookHandler_DeleteRequest_Failure(t *testing.T) {
	mockClient := &OverseerClientMock{}
	mockClient.On("DeclineRequest", mock.AnythingOfType("string")).Return(nil)
	mockClient.On("DeleteRequest", mock.AnythingOfType("string")).Return(errors.New("Error while deleting request"))
	handler := &Handler{
		OverseerrClient: mockClient,
		BannedMediaIDs:  []string{"999"},
		DeleteRequests:  true,
	}
	payload := `{"request_id": "1", "tmdbid": "999", "tvdbid": "888"}`
	req := httptest.NewRequest("POST", "/", bytes.NewBuffer([]byte(payload)))
	rec := httptest.NewRecorder()

	handler.WebhookHandler(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	mockClient.AssertCalled(t, "DeclineRequest", "1")
	mockClient.AssertNotCalled(t, "DeleteRequest", "1")
	mockClient.AssertExpectations(t)
}
