package lib

import (
	"testing"
	"net/http"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"github.com/gorilla/websocket"
	"fmt"
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
	roomHandler := NewRoom()
	server := httptest.NewServer(roomHandler)
	dialer := websocket.Dialer{}
	message := "This is a test message"
	url := fmt.Sprintf("ws://%s/room", server.Listener.Addr().String())
	connection, response, error := dialer.Dial(url, nil)
	if error != nil {
		t.Fatal(error)
	}
	//assert that we established a full duplex connection
	assert.Equal(t, response.Status,"101 Switching Protocols" )
	connection.WriteMessage(websocket.TextMessage, []byte(message))
	messageFromChannel := <-roomHandler.broadcast
	//assert that we have message recieved
	assert.Equal(t, message, string(messageFromChannel))
}