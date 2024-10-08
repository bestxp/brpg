package main

import (
	"net/http"
	"time"

	"github.com/bestxp/brpg/internal/game"
	"github.com/bestxp/brpg/internal/infra/network"
	engine "github.com/bestxp/brpg/pkg"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// // Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	id   string
	hub  *Hub
	conn *network.Network
	send chan *engine.Event
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump(world *game.World) {
	defer func() {
		event := &engine.Event{
			Type: engine.Event_type_exit,
			Data: &engine.Event_Exit{
				Exit: &engine.EventExit{PlayerId: c.id},
			},
		}
		world.HandleEvent(event)
		c.hub.broadcast <- event

		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.Conn.SetReadLimit(maxMessageSize)
	c.conn.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.Conn.SetPongHandler(func(string) error { c.conn.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, event, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		c.hub.broadcast <- event // ?
		world.HandleEvent(event)
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
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
				// The hub closed the channel.
				c.conn.CloseMessage()
				return
			}

			w, err := c.conn.BathBinary()
			if err != nil {
				return
			}
			w.Send(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Send(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.PingMessage(); err != nil {
				return
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
func serveWs(hub *Hub, world *game.World, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error().Err(err)
		return
	}
	net := network.NewNetwork(conn)

	id := world.AddPlayer()
	client := &Client{id: id, hub: hub, conn: net, send: make(chan *engine.Event, 256)}
	client.hub.register <- client

	event := &engine.Event{
		Type: engine.Event_type_init,
		Data: &engine.Event_Init{
			Init: &engine.EventInit{
				PlayerId: id,
				Units:    world.Units,
			},
		},
	}
	err = net.Send(event)
	if err != nil {
		//todo: remove unit
		log.Error().Err(err)
	}

	unit := world.Units[id]
	event = &engine.Event{
		Type: engine.Event_type_connect,
		Data: &engine.Event_Connect{
			Connect: &engine.EventConnect{Unit: unit},
		},
	}
	hub.broadcast <- event

	// Allow collection of memory referenced by the caller by doing all work
	// in new goroutines.
	go client.writePump()
	go client.readPump(world)
}
