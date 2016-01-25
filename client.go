package main

import (
	"fmt"

	r "github.com/dancannon/gorethink"
	"github.com/gorilla/websocket"
)

// FindHandler findHandler object
type FindHandler func(string) (Handler, bool)

// Message is a message
type Message struct {
	Name string      `json:"name"`
	Data interface{} `json:"data"`
}

// Client is a client
type Client struct {
	send        chan Message
	socket      *websocket.Conn
	findHandler FindHandler
	session     *r.Session
}

func (client *Client) Read() {
	var message Message
	for {
		if err := client.socket.ReadJSON(&message); err != nil {
			// fmt.Println(err)
			// fmt.Println("breaking out of client.Read loop")
			break
		}
	}
	if handler, found := client.findHandler(message.Name); found {
		fmt.Println("found: ", found, " message.Name: ", message.Name)
		handler(client, message.Data)
	}
	// fmt.Println("Closing client.socket")
	client.socket.Close()
}

func (client *Client) Write() {
	for msg := range client.send {
		if err := client.socket.WriteJSON(msg); err != nil {
			// fmt.Println(err)
			// fmt.Println("breaking out of client.Write loop")
			break
		}
	}
	// fmt.Println("Closing client.socket")
	client.socket.Close()
}

// NewClient creates a new Client
func NewClient(socket *websocket.Conn, findHandler FindHandler, session *r.Session) *Client {
	return &Client{
		send:        make(chan Message),
		socket:      socket,
		findHandler: findHandler,
		session:     session,
	}
}
