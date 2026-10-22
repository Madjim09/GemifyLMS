package handler_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestHandleWebhookEvent_NewUser(t *testing.T) {
	deps := newTestDeps(t)

	r := gin.New()
	r.POST("/webhook/event", deps.handler.HandleWebhookEvent)

	body := `{
		"actor": {"mbox":"mailto:webhook@example.com","name":"Webhook User"},
		"verb":  {"id":"http://adlnet.gov/expapi/verbs/passed","display":{"en":"passed"}},
		"object":{"id":"module-1","objectType":"Activity"}
	}`
	req := httptest.NewRequest(http.MethodPost, "/webhook/event", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d; body: %s", w.Code, http.StatusOK, w.Body.String())
	}
}

func TestHandleWebhookEvent_WithAccount(t *testing.T) {
	deps := newTestDeps(t)

	r := gin.New()
	r.POST("/webhook/event", deps.handler.HandleWebhookEvent)

	body := `{
		"actor": {"name":"Account User","account":{"homePage":"https://example.com","name":"acc-user-123"}},
		"verb":  {"id":"http://adlnet.gov/expapi/verbs/completed","display":{"en":"completed"}},
		"object":{"id":"course-1","objectType":"Activity"}
	}`
	req := httptest.NewRequest(http.MethodPost, "/webhook/event", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d; body: %s", w.Code, http.StatusOK, w.Body.String())
	}
}

func TestHandleWebhookEvent_BadJSON(t *testing.T) {
	deps := newTestDeps(t)

	r := gin.New()
	r.POST("/webhook/event", deps.handler.HandleWebhookEvent)

	req := httptest.NewRequest(http.MethodPost, "/webhook/event", strings.NewReader(`not json`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestHandleWebhookEvent_NoActor(t *testing.T) {
	deps := newTestDeps(t)

	r := gin.New()
	r.POST("/webhook/event", deps.handler.HandleWebhookEvent)

	body := `{
		"actor": {},
		"verb":  {"id":"http://adlnet.gov/expapi/verbs/passed"},
		"object":{"id":"module-1"}
	}`
	req := httptest.NewRequest(http.MethodPost, "/webhook/event", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d; body: %s", w.Code, http.StatusBadRequest, w.Body.String())
	}
}

func TestHandleWebhookEvent_UnknownVerb(t *testing.T) {
	deps := newTestDeps(t)

	r := gin.New()
	r.POST("/webhook/event", deps.handler.HandleWebhookEvent)

	body := `{
		"actor": {"mbox":"mailto:unknown@example.com","name":"Unknown Verb User"},
		"verb":  {"id":"http://example.com/custom/verb","display":{"en":"custom"}},
		"object":{"id":"item-1","objectType":"Activity"}
	}`
	req := httptest.NewRequest(http.MethodPost, "/webhook/event", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Unknown verb is treated as a custom event type — should still succeed
	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d; body: %s", w.Code, http.StatusOK, w.Body.String())
	}
}

func TestHandleWebhookEvent_ExistingUserSecondEvent(t *testing.T) {
	deps := newTestDeps(t)

	r := gin.New()
	r.POST("/webhook/event", deps.handler.HandleWebhookEvent)

	makeBody := func(verb string) string {
		return `{"actor":{"mbox":"mailto:repeat@example.com","name":"Repeat User"},"verb":{"id":"` + verb + `"},"object":{"id":"m1"}}`
	}

	// First event creates user
	req1 := httptest.NewRequest(http.MethodPost, "/webhook/event", strings.NewReader(makeBody("http://adlnet.gov/expapi/verbs/passed")))
	req1.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(httptest.NewRecorder(), req1)

	// Second event should find existing user
	req2 := httptest.NewRequest(http.MethodPost, "/webhook/event", strings.NewReader(makeBody("http://adlnet.gov/expapi/verbs/completed")))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	if w2.Code != http.StatusOK {
		t.Errorf("status = %d, want %d; body: %s", w2.Code, http.StatusOK, w2.Body.String())
	}
}
