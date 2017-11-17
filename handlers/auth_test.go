package handlers

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockHandler struct{}

func (m *MockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
func TestMustAuth(t *testing.T) {
	req, err := http.NewRequest("GET", "/chat", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	roomHandler := NewRoom()
	ah := authHandler{next: roomHandler}
	handler := MustAuth(http.Handler(&ah))
	handler.ServeHTTP(rr, req)
	assert.Equal(t, rr.Code, http.StatusTemporaryRedirect)
}
func TestAuthHandlerServeHTTPWithCookie(t *testing.T) {
	rr := httptest.NewRecorder()
	http.SetCookie(rr, &http.Cookie{Name: "auth", Value: "test"})
	//request with cookie
	request := &http.Request{Header: http.Header{"Cookie": rr.HeaderMap["Set-Cookie"]}}
	mh := &MockHandler{}
	ah := authHandler{next: mh}
	handler := http.Handler(&ah)
	handler.ServeHTTP(rr, request)
	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, rr.Header().Get("Set-Cookie"), "auth=test")
}

func TestAuthHandlerServeHTTPWithOutCookie(t *testing.T) {
	rr := httptest.NewRecorder()
	//request without cookie
	request := &http.Request{}
	mh := &MockHandler{}
	ah := authHandler{next: mh}
	handler := http.Handler(&ah)
	handler.ServeHTTP(rr, request)
	assert.Equal(t, rr.Code, http.StatusTemporaryRedirect)
	assert.Equal(t, rr.Header().Get("Set-Cookie"), "")
	assert.Equal(t, rr.Header().Get("Location"), "/login")
}

func TestCheckHttpErrorWithOutCookie(t *testing.T) {
	rr := httptest.NewRecorder()
	request := &http.Request{}
	httpError := CheckCookie(rr, request)
	assert.Equal(t, httpError, http.StatusTemporaryRedirect)
	assert.Equal(t, rr.Header().Get("Location"), "/login")
}

func TestCheckHttpErrorWithCookie(t *testing.T) {
	rr := httptest.NewRecorder()
	http.SetCookie(rr, &http.Cookie{Name: "auth", Value: "test"})
	//request with cookie
	request := &http.Request{Header: http.Header{"Cookie": rr.HeaderMap["Set-Cookie"]}}
	httpError := CheckCookie(rr, request)
	assert.Equal(t, httpError, 0)
	assert.Equal(t, rr.Header().Get("Location"), "")
}
func TestCheckHttpErrorWithWrongCookie(t *testing.T) {
	rr := httptest.NewRecorder()
	http.SetCookie(rr, &http.Cookie{Name: "wrongcookie", Value: "test"})
	//request with cookie
	request := &http.Request{Header: http.Header{"Cookie": rr.HeaderMap["Set-Cookie"]}}
	httpError := CheckCookie(rr, request)
	assert.Equal(t, httpError, http.StatusTemporaryRedirect)
	assert.Equal(t, rr.Header().Get("Location"), "/login")
}
