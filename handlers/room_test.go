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
	room := NewRoom(false)
	assert.IsType(t, &Room{}, room)
}
func TestRoomServeHTTPWithHTTPRequest(t *testing.T) {
	req, err := http.NewRequest("GET", "/room", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	roomHandler := NewRoom(false)
	roomHandler.ServeHTTP(rr, req)
	assert.Equal(t, rr.Code, http.StatusBadRequest)
}
func TestRoomCreateClient(t *testing.T) {
	req, err := http.NewRequest("GET", "/room", nil)
	Header := map[string][]string{
				"Connection": {"upgrade"},
				"Upgrade": {"websocket"},
				"Sec-Websocket-Version": {"13"},
				"Sec-Websocket-Key": {"key"},
			}

	req.Header = Header
	if err != nil {
		t.Fatal(err)
	}
	rw := httptest.NewRecorder()
	roomHandler := NewRoom(false)
	_, err = roomHandler.createClient(rw, req)
	assert.Error(t, err)
	assert.Equal(t, len(roomHandler.clients), 0)
}
func TestRoomCreateClientWithWebSocket(t *testing.T) {
	roomHandler := NewRoom(true)
	server := httptest.NewServer(roomHandler)
	dialer := websocket.Dialer{}
	urlString := fmt.Sprintf("ws://%s/room", server.Listener.Addr().String())
	_, _ , error := dialer.Dial(urlString, nil)
	assert.NoError(t, error)
	assert.Equal(t, len(roomHandler.clients), 1)
	fmt.Println(len(roomHandler.clients))
}
func TestRoomServeHTTPWithWebSocketRequest(t *testing.T) {
	var cookies []*http.Cookie
	roomHandler := NewRoom(true)
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
	assert.Equal(t, len(roomHandler.clients), 1)
}
//func TestRoomServeHTTPWithWebSocketRequestWithInvalidCookie(t *testing.T) {
//	var cookies []*http.Cookie
//	roomHandler := NewRoom(false)
//	server := httptest.NewServer(roomHandler)
//	dialer := websocket.Dialer{}
//	urlString := fmt.Sprintf("ws://%s/room", server.Listener.Addr().String())
//
//	cookieJar, _ := cookiejar.New(nil)
//	cookie := http.Cookie{Name: "invalid", Value: "test"}
//	cookies = append(cookies, &cookie)
//	urlParsed, _ := url.Parse(urlString)
//
//	cookieJar.SetCookies(urlParsed, cookies)
//	dialer.Jar = cookieJar
//
//	connection, _, error := dialer.Dial(urlString, nil)
//	if error != nil {
//		t.Fatalf("Error encountered %s", error)
//	}
//	connection.WriteMessage(websocket.TextMessage, []byte("my message"))
//	connection.Close()
//	messageFromChannel := <-roomHandler.broadcast
//	close(roomHandler.broadcast)
//	fmt.Println(messageFromChannel)
//	assert.Error(t, error)
//}