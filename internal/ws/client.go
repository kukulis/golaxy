// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ws

import (
	"bytes"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// client is a middleman between the websocket connection and the hub.
type client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	token string
}

func (c *client) readPump() {
	defer func() {
		c.hub.dispatcher.Dispatch(NewWsEvent(EventClientDisconnected, c.token, nil))
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			c.hub.dispatcher.Dispatch(NewWsEvent(EventMessageArrivedError, c.token, message))
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		// TODO consider trim and replace logic through dispatcher?
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		checkSendEvent := NewWsCheckSendEvent(EventMessageArrived, c.token, message)
		c.hub.dispatcher.Dispatch(checkSendEvent)

		if checkSendEvent.Send {
			c.hub.broadcast <- message
			c.hub.dispatcher.Dispatch(NewWsEvent(EventMessageSentToHub, c.token, message))
		}
	}
}

func (c *client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.hub.dispatcher.Dispatch(NewWsEvent(EventSendMessageChannelFail, c.token, message))
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.hub.dispatcher.Dispatch(NewWsEvent(EventBeforeMessageSend, c.token, message))

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				c.hub.dispatcher.Dispatch(NewWsEvent(EventWriterInitializeError, c.token, message))
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				c.hub.dispatcher.Dispatch(NewWsEvent(EventMessageSendError, c.token, message))
				return
			}
			c.hub.dispatcher.Dispatch(NewWsEvent(EventMessageSendSuccess, c.token, message))
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
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
