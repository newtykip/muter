package main

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/kardianos/service"
)

var upgrader = websocket.Upgrader{}

func handleWebsocket(w http.ResponseWriter, r *http.Request, logger service.Logger) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			logger.Error("read:", err)
			break
		}
		logger.Infof("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			logger.Error("write:", err)
			break
		}
	}
}
