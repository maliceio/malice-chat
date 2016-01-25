package main

import (
	"fmt"
	"net/http"

	r "github.com/dancannon/gorethink"
	"github.com/gorilla/websocket"
)

// Handler is a handler
type Handler func(*Client, interface{})

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// Router is router
type Router struct {
	rules   map[string]Handler
	session *r.Session
}

func handleError(socket *websocket.Conn, err error) {
	var outMessage Message
	if err != nil {
		outMessage = Message{Name: "error", Data: err}
		assert(socket.WriteJSON(outMessage))
	}
}

// NewRouter creates a new Router
func NewRouter(session *r.Session) *Router {
	return &Router{
		rules:   make(map[string]Handler),
		session: session,
	}
}

// Handle handles thangs
func (rout *Router) Handle(msgName string, handler Handler) {
	rout.rules[msgName] = handler
}

// FindHandler finds handler
func (rout *Router) FindHandler(msgName string) (Handler, bool) {
	handler, found := rout.rules[msgName]
	return handler, found
}

// ServeHTTP serves it up fresh
func (rout *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}
	client := NewClient(socket, rout.FindHandler, rout.session)
	go client.Write()
	client.Read()
}
