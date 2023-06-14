package main

import "github.com/gorilla/websocket"

//client is a single chatting user

type client struct {
	socket *websocket.Conn

	// channel for messages
	recieve chan []byte

	//the room the client is chatting in
	room *room
}

type room struct {
}

func (c *client) read() {
	defer c.socket.Close()
	for {
		_, msg, err := c.socket.ReadMessage()
		if err != nil {
			return
		}
		c.room.forward <- msg
	}
}

func (c *client) write() {
	defer c.socket.Close()

	for msg := range c.recieve {
		err := c.socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return
		}
	}
}
