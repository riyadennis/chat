package handlers

import (
	//"fmt"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	//"net/url"
	"fmt"
	"net/url"
	"net/http/cookiejar"
)

func TestNewRoom(t *testing.T) {
	room := NewRoom()
	assert.IsType(t, &Room{}, room)
}
func TestRoomServeHTTPWithHTTPRequest(t *testing.T) {
	req, err := http.NewRequest("GET", "/room", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	roomHandler := NewRoom()
	roomHandler.ServeHTTP(rr, req)
	assert.Equal(t, rr.Code, http.StatusBadRequest)
}
func TestRoomServeHTTPWithWebSocketRequest(t *testing.T) {
	var cookies []*http.Cookie
	roomHandler := NewRoom()
	server := httptest.NewServer(roomHandler)
	dialer := websocket.Dialer{}
	urlString := fmt.Sprintf("ws://%s/room", server.Listener.Addr().String())

	cookieJar, _ := cookiejar.New(nil)
	cookie := http.Cookie{Name: "auth", Value: "test"}
	cookies = append(cookies, &cookie)
	urlParsed, _ := url.Parse(urlString)

	cookieJar.SetCookies(urlParsed, cookies)
	dialer.Jar = cookieJar

	_, response, error := dialer.Dial(urlString, nil)
	if error != nil {
		t.Fatalf("Error encountered %s", error)
	}
	//assert that we established a full duplex connection
	assert.Equal(t, response.Status, "101 Switching Protocols")
}