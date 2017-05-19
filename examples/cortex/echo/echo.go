package echo

import (
	"context"
	"github.com/owulveryck/cortical"
	"log"
)

// Echo is a dummy type that reads a message, wait for some time and sends ret back
type Echo struct {
	pong string
	c    chan []byte
}

// NewCortex is filling the  ...
func NewCortex(ctx context.Context) (cortical.GetInfoFromCortexFunc, cortical.SendInfoToCortex) {
	c := make(chan []byte)
	echo := &Echo{
		pong: "pong",
		c:    c,
	}
	return echo.Get, echo.Receive
}

// Get ...
func (e *Echo) Get(ctx context.Context) chan []byte {
	return e.c
}

// Receive ...
func (e *Echo) Receive(ctx context.Context, b *[]byte) {
	log.Printf("[%v] received (%v)", e.pong, ctx)
	e.c <- []byte(e.pong)
}
