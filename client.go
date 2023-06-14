package main

import (
	"encoding/json"

	"github.com/gorilla/websocket"
)

//client is a single chatting user

type client struct {
	socket *websocket.Conn

	// channel for messages
	recieve chan []byte

	//the room the client is chatting in
	room *room

	name string
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

type fullMsg struct {
	Msg  string `json:"msg"`
	Name string `json:"name"`
}

func (c *client) write() {
	defer c.socket.Close()

	for msg := range c.recieve {
		message := &fullMsg{
			Msg:  string(msg),
			Name: c.name,
		}
		jsonData, err := json.Marshal(message)
		if err != nil {
			return
		}
		err = c.socket.WriteMessage(websocket.TextMessage, jsonData)
		if err != nil {
			return
		}
	}
}
