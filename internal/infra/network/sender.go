package network

import (
	"fmt"
	"io"
	"time"

	"github.com/bestxp/brpg/pkg"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
)

type Network struct {
	Conn *websocket.Conn
}

type bathSender struct {
	w io.WriteCloser
}

func (b *bathSender) Send(e *pkg.Event) error {
	message, err := proto.Marshal(e)
	if err != nil {
		return fmt.Errorf("failed encode event: %w", err)
	}

	if _, err = b.w.Write(message); err != nil {
		return fmt.Errorf("failed sent message: %w", err)
	}

	return nil
}

func (b *bathSender) Close() error {
	return b.w.Close()
}

type MessageSender interface {
	Send(event *pkg.Event) error
	Close() error
}

func FromHost(host string) *Network {
	for i := 0; i < 10; i++ {
		c, _, err := websocket.DefaultDialer.Dial("ws://"+host+":3000/ws", nil)
		if err != nil {
			time.Sleep(time.Second * time.Duration(i))
			log.Debug().Msgf("Connection timeout, retry %s", err.Error())
		} else {
			return NewNetwork(c)
		}
	}
	return nil
}

func NewNetwork(c *websocket.Conn) *Network {
	return &Network{Conn: c}
}

func (n *Network) Send(e *pkg.Event) error {
	message, err := proto.Marshal(e)
	if err != nil {
		return fmt.Errorf("failed encode event: %w", err)
	}

	if err = n.Conn.WriteMessage(websocket.BinaryMessage, message); err != nil {
		return fmt.Errorf("failed sent message: %w", err)
	}

	return nil
}

func (n *Network) Close() error {
	return n.Conn.Close()
}

func (n *Network) ReadMessage() (int, *pkg.Event, error) {
	in, message, err := n.Conn.ReadMessage()
	if err != nil {
		return 0, nil, err
	}

	var event pkg.Event
	err = proto.Unmarshal(message, &event)
	if err != nil {
		log.Error().Err(err)
	}
	return in, &event, nil
}

func (n *Network) CloseMessage() error {
	return n.Conn.WriteMessage(websocket.CloseMessage, []byte{})
}

func (n *Network) PingMessage() error {
	return n.Conn.WriteMessage(websocket.PingMessage, []byte{})
}

func (n *Network) BathBinary() (MessageSender, error) {
	w, err := n.Conn.NextWriter(websocket.BinaryMessage)
	if err != nil {
		return nil, err
	}
	return &bathSender{w: w}, nil
}

func (n *Network) SetWriteDeadline(add time.Time) {
	n.Conn.SetWriteDeadline(add)
}
