package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/riyadennis/chat/config"
	"github.com/stretchr/testify/assert"
)

func TestNewLoginProviderHandler(t *testing.T) {
	config := &config.Config{}
	lph := NewLoginHandler(config)
	assert.IsType(t, &loginHandler{}, lph)
}
func TestLoginHandlerServeHTTPWillGiveOKForValidLoginURL(t *testing.T) {
	req, err := http.NewRequest("GET", "/auth/login/google", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	config := &config.Config{}
	lh := NewLoginHandler(config)
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

	config := &config.Config{}
	lph := NewLoginHandler(config)

	handler := MustAuth(http.Handler(lph))
	handler.ServeHTTP(rr, req)
	assert.Equal(t, rr.Code, http.StatusTemporaryRedirect)
}
