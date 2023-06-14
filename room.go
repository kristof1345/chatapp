package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

type room struct {
	//holds current clients in the room
	clients map[*client]bool

	//channel for clients who want to join
	join chan *client

	//leave channel
	leave chan *client

	//forward message
	forward chan []byte
}

func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		leave:   make(chan *client),
		join:    make(chan *client),
		clients: make(map[*client]bool),
	}
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			r.clients[client] = true
		case client := <-r.leave:
			delete(r.clients, client)
			close(client.recieve)
		case msg := <-r.forward:
			for client := range r.clients {
				client.recieve <- msg
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: messageBufferSize}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP: ", err)
		return
	}
	client := &client{
		socket:  socket,
		recieve: make(chan []byte, messageBufferSize),
		room:    r,
		name:    "client" + strconv.Itoa(len(r.clients)),
	}

	r.join <- client
	defer func() { r.leave <- client }() // if client exits disconnect him

	go client.write()
	client.read()
}
