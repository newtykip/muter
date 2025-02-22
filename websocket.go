package main

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/kardianos/service"
)

var upgrader = websocket.Upgrader{}

func handleWebsocket(w http.ResponseWriter, r *http.Request, state chan bool, logger service.Logger) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error("upgrade:", err)
		return
	}
	defer c.Close()
	for muted := range state {
		if err := c.WriteJSON(muted); err != nil {
			logger.Error("write:", err)
			return
		}
	}
}
