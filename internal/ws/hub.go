// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ws

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"glaktika.eu/galaktika/pkg/util"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *client

	// Unregister requests from clients.
	unregister chan *client

	dispatcher *util.Dispatcher
}

func NewHub(dispatcher *util.Dispatcher) *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *client),
		unregister: make(chan *client),
		clients:    make(map[*client]bool),
		dispatcher: dispatcher,
	}
}

// Run this makes registering, unregistering and message sending actions never happen in parallel
func (h *Hub) Run() {
	for {
		select {
		case c := <-h.register:
			h.clients[c] = true
		case c := <-h.unregister:
			if _, ok := h.clients[c]; ok {
				delete(h.clients, c)
				close(c.send)
			}
		case message := <-h.broadcast:

			checkReceiverEvent := NewWsCheckSendEvent(EventCheckReceiver, "", message)
			h.dispatcher.Dispatch(checkReceiverEvent)
			for c := range h.clients {
				if checkReceiverEvent.ReceiverToken != "" && checkReceiverEvent.ReceiverToken != c.token {
					continue
				}
				select {
				case c.send <- message:
				default:
					close(c.send)
					delete(h.clients, c)
				}
			}
		}
	}
}

// RegisterWebSocketClient handles websocket requests from the peer.
func RegisterWebSocketClient(hub *Hub, w http.ResponseWriter, r *http.Request, token string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	c := &client{hub: hub, conn: conn, send: make(chan []byte, 256), token: token}
	c.hub.register <- c

	go c.writePump()
	go c.readPump()
}
