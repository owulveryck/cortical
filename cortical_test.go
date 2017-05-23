package cortical

import (
	"github.com/gorilla/websocket"
	"time"

	"context"
	"github.com/gorilla/mux"
	"net/http/httptest"
	"net/url"
	"testing"
)

var (
	testServer *httptest.Server
	tsURL      *url.URL
)

// Echo is a dummy type that reads a message, wait for some time and sends ret back
type Echo struct {
	c chan []byte
}

func new() *Echo {
	c := make(chan []byte)
	return &Echo{
		c: c,
	}
}

// NewCortex is filling the  ...
func (e *Echo) NewCortex(ctx context.Context) (GetInfoFromCortexFunc, SendInfoToCortex) {
	return e.Get, e.Receive
}

// Get ...
func (e *Echo) Get(ctx context.Context) chan []byte {
	return e.c
}

// Receive ...
func (e *Echo) Receive(ctx context.Context, b *[]byte) {
	e.c <- *b
}
func init() {
	router := mux.NewRouter().StrictSlash(true)
	brain := &Cortical{
		Upgrader: websocket.Upgrader{},
		Cortexes: []Cortex{new()},
	}

	router.
		Methods("GET").
		Path("/ws").
		Name("WebSocket").
		HandlerFunc(brain.ServeWS)

	testServer = httptest.NewServer(router) //Creating new server with the user handlers
	tsURL, _ = url.Parse(testServer.URL)
}

func TestPingPong(t *testing.T) {
	wsURL := url.URL{Scheme: "ws", Host: tsURL.Host, Path: "/ws"}
	c, _, err := websocket.DefaultDialer.Dial(wsURL.String(), nil)
	if err != nil {
		t.Fatalf("Cannot connect to the websocket %v", err)

	}
	defer c.Close()
	if err := c.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(2*time.Second)); err != nil {
		t.Errorf("write close: %v", err)
	}
	err = c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseGoingAway, ""))
	if err != nil {
		t.Errorf("write close: %v", err)
	}
}
func TestClosing(t *testing.T) {
	// Now test the websocket
	wsURL := url.URL{Scheme: "ws", Host: tsURL.Host, Path: "/ws"}
	c, _, err := websocket.DefaultDialer.Dial(wsURL.String(), nil)
	if err != nil {
		t.Errorf("Cannot connect to the websocket %v", err)
	}
	defer c.Close()
}

func TestServeWS(t *testing.T) {
	// Now test the websocket
	test := []byte("test")
	wsURL := url.URL{Scheme: "ws", Host: tsURL.Host, Path: "/ws"}
	c, _, err := websocket.DefaultDialer.Dial(wsURL.String(), nil)
	if err != nil {
		t.Errorf("Cannot connect to the websocket %v", err)
	}
	defer c.Close()
	done := make(chan bool)

	go func() {
		defer close(done)
		tm, message, err := c.ReadMessage()
		if err != nil {
			t.Fatal("Error in the message reception: %v (type %v)", err, tm)
		}
		if string(message) != string(test) {
			t.Fatal("Message received should be the same as the message sent")
		}
		done <- true
	}()

	err = c.WriteMessage(websocket.TextMessage, test)
	if err != nil {
		t.Errorf("Cannot write %v to websocket: %v", test, err)
	}

	<-done
	err = c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseGoingAway, ""))
	if err != nil {
		t.Errorf("write close: %v", err)
	}
}
