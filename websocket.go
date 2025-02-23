package main

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/kardianos/service"
	"github.com/moutend/go-wca/pkg/wca"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleWebsocket(w http.ResponseWriter, r *http.Request, aev *wca.IAudioEndpointVolume, logger service.Logger) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error("upgrade:", err)
		return
	}
	defer c.Close()

	// send initial state
	state.Lock()
	c.WriteJSON(state.muted)
	state.Unlock()

	// subscribe to state changes
	ch := state.subscribe()
	defer state.unsubscribe(ch)

	// done channel for cleanup
	done := make(chan struct{})
	defer close(done)

	// handle incoming messages
	go func() {
		for {
			var muted bool
			if err := c.ReadJSON(&muted); err != nil {
				logger.Error("read:", err)
				done <- struct{}{}
				return
			}

			if err := aev.SetMute(muted, nil); err != nil {
				logger.Error("SetMute:", err)
				done <- struct{}{}
				return
			}
		}
	}()

	// handle state updates
	for {
		select {
		case muted := <-ch:
			if err := c.WriteJSON(muted); err != nil {
				logger.Error("write:", err)
				return
			}
		case <-done:
			return
		}
	}
}
