package handler

import (
	"bytes"
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
	return nil
}

func (m *OverseerClientMock) DeleteRequest(requestID string) error {
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
	mockClient := new(OverseerClientMock)
	mockClient.On("DeclineRequest", "1").Return(nil)

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

}

func TestWebhookHandler_DeclineRequest_Failure(t *testing.T) {

}

func TestWebhookHandler_DeleteRequest_Success(t *testing.T) {

}

func TestWebhookHandler_DeleteRequest_Failure(t *testing.T) {

}
