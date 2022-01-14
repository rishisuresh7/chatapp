package chat

import (
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn    *websocket.Conn
	channel chan *Message
	userId  int64
	hub     *ChatHub
}

func NewChatClient(h *ChatHub, ws *websocket.Conn, userId int64) *Client {
	return &Client{
		conn:    ws,
		channel: make(chan *Message),
		userId:  userId,
		hub:     h,
	}
}

func (c *Client) GetId() int64 {
	return c.userId
}

func (c *Client) Reader() {
	defer func() {
		c.hub.DeRegister(c)
		c.conn.Close()
	}()
	for {
		var msg Message
		err := c.conn.ReadJSON(&msg)
		if err != nil {
			return
		}

		for client := range c.hub.clients {
			if client.userId == msg.To {
				msg.From = c.userId
				msg.Time = time.Now().Unix()
				client.channel <- &msg
			}
		}
	}
}

func (c *Client) Writer() {
	defer c.conn.Close()
	for {
		select {
		case msg := <-c.channel:
			if msg != nil {
				err := c.conn.WriteJSON(msg)
				if err != nil {
					return
				}
			}
		}
	}
}

type Message struct {
	From     int64  `json:"sender"`
	To       int64	`json:"receiver"`
	Time     int64 `json:"time"`
	Message  string `json:"message"`
}
