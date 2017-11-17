package handlers

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewLoginProviderHandler(t *testing.T) {
	lph := NewLoginHandler()
	assert.IsType(t, &loginHandler{}, lph)
}
func TestLoginHandlerServeHTTPWillGiveOKForValidLoginURL(t *testing.T) {
	req, err := http.NewRequest("GET", "/auth/login/google", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	lh := NewLoginHandler()
	handler := MustAuth(http.Handler(lh))
	handler.ServeHTTP(rr, req)
	assert.Equal(t, rr.Code, http.StatusTemporaryRedirect)
}
func TestLoginHandlerServeHTTPWillGiveErrorForInvalidValidLoginURL(t *testing.T) {
	req, err := http.NewRequest("GET", "/auth/userlogin/facebook", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	lph := NewLoginHandler()

	handler := MustAuth(http.Handler(lph))
	handler.ServeHTTP(rr, req)
	assert.Equal(t, rr.Code, http.StatusTemporaryRedirect)
}
