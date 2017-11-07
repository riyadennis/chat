package handlers

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"net/http"
	"fmt"
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
func TestGetLoginURL(t *testing.T) {
	url, err := getLoginURL("sdsfd")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(url)
}
