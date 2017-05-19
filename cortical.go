// Package cortical is a utility that handles websocket connexions and dispatch
// all the []byte received to the consumers
// It also get all the informations of the producers and sends them back to the websocket
//
// Cortex is the interface that every "processor" must implement.
// The first output is a send function and the second argument is a receive function
// Every Cortex will receive every elements through their receive function
// Every send function will be called in an endless loop.
//
// The context passed to the Cortex holds a uuid encoded into a string that is specific for each connection
// The UUID can be retrieved by calling
//
//   uuid := ctx.Value(ContextKeyType(ContextKey))
package cortical

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"net/http"
)

// Cortical specifies how to upgrade an HTTP connection to a Websocket connection
// as well as the action to be performed on receive a []byte
type Cortical struct {
	Upgrader websocket.Upgrader
	Cortexs  []func(context.Context) (GetInfoFromCortexFunc, SendInfoToCortex)
}

// ContextKeyType is the type of the key of the context
type ContextKeyType string

// ContextKey is the key name where the session is stored
const ContextKey = "uuid"

// GetInfoFromCortexFunc is the method implenented by a chatter to send objects
type GetInfoFromCortexFunc func(ctx context.Context) chan []byte

// SendInfoToCortex is the method implemented by a chatter to receive objects
type SendInfoToCortex func(context.Context, *[]byte)

type httpErr struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
}

func handleErr(w http.ResponseWriter, err error, status int) {
	msg, err := json.Marshal(&httpErr{
		Msg:  err.Error(),
		Code: status,
	})
	if err != nil {
		msg = []byte(err.Error())
	}
	http.Error(w, string(msg), status)
}

// ServeWS is the dispacher function
func (wsd *Cortical) ServeWS(w http.ResponseWriter, r *http.Request) {
	ctx := context.WithValue(r.Context(), ContextKeyType(ContextKey), uuid.New().String())
	conn, err := wsd.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		handleErr(w, err, http.StatusInternalServerError)
		return
	}
	defer conn.Close()
	var senders = make([]GetInfoFromCortexFunc, 0)
	var receivers = make([]SendInfoToCortex, 0)
	for _, cortex := range wsd.Cortexs {
		snd, rcv := cortex(ctx)
		if snd != nil {
			senders = append(senders, snd)
		}
		if rcv != nil {
			receivers = append(receivers, rcv)
		}
	}
	rcvsNum := len(receivers)
	sndrsNum := len(senders)
	var stop []chan struct{}
	for i := 0; i < sndrsNum+rcvsNum; i++ {
		s := make(chan struct{})
		stop = append(stop, s)
	}
	rcv := make(chan []byte, 1)
	sendersChan := make([]<-chan []byte, sndrsNum)
	chans := fanOut(rcv, rcvsNum, 1)
	for i := 0; i < sndrsNum; i++ {
		sendersChan[i] = senders[i](ctx)
	}
	for i := range chans {
		receive(ctx, chans[i], stop[i+sndrsNum], receivers[i])
	}
	done := make(chan struct{}, 1)
	send := merge(done, sendersChan...)
	closed := make(chan struct{}, 2)
	go func() {
		for {
			p := <-send
			err := conn.WriteMessage(websocket.TextMessage, p)
			if err != nil {
				if websocket.IsCloseError(err,
					websocket.CloseNormalClosure,
					websocket.CloseGoingAway,
					websocket.CloseNormalClosure,
					websocket.CloseGoingAway,
					websocket.CloseProtocolError,
					websocket.CloseUnsupportedData,
					websocket.CloseNoStatusReceived,
					websocket.CloseAbnormalClosure,
					websocket.CloseInvalidFramePayloadData,
					websocket.ClosePolicyViolation,
					websocket.CloseMessageTooBig,
					websocket.CloseMandatoryExtension,
					websocket.CloseInternalServerErr,
					websocket.CloseServiceRestart,
					websocket.CloseTryAgainLater,
					websocket.CloseTLSHandshake,
					websocket.CloseNoStatusReceived) {
					closed <- struct{}{}
					return
				}
				if err == websocket.ErrCloseSent {
					closed <- struct{}{}
					return
				}
				// Temporary failure, nevermind
				continue
			}
		}
	}()
	go func() {
		for {
			MessageType, p, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err,
					websocket.CloseNormalClosure,
					websocket.CloseGoingAway,
					websocket.CloseNormalClosure,
					websocket.CloseGoingAway,
					websocket.CloseProtocolError,
					websocket.CloseUnsupportedData,
					websocket.CloseNoStatusReceived,
					websocket.CloseAbnormalClosure,
					websocket.CloseInvalidFramePayloadData,
					websocket.ClosePolicyViolation,
					websocket.CloseMessageTooBig,
					websocket.CloseMandatoryExtension,
					websocket.CloseInternalServerErr,
					websocket.CloseServiceRestart,
					websocket.CloseTryAgainLater,
					websocket.CloseTLSHandshake,
					websocket.CloseNoStatusReceived) {
					closed <- struct{}{}
					return
				}
				if err == websocket.ErrCloseSent {
					closed <- struct{}{}
					return
				}
				// Temporary failure, nevermind
				continue
			}
			if MessageType != websocket.TextMessage {
				handleErr(w, errors.New("Only text []byte are supported"), http.StatusNotImplemented)
				continue
			}
			rcv <- p
		}
	}()
	<-closed
	done <- struct{}{}
	for i := 0; i < sndrsNum+rcvsNum; i++ {
		stop[i] <- struct{}{}
	}
}

func receive(ctx context.Context, msg <-chan []byte, stop chan struct{}, f SendInfoToCortex) {
	go func() {
		for {
			select {
			case b := <-msg:
				f(ctx, &b)
			case <-stop:
				return
			}
		}
	}()
}
