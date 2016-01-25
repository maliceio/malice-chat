package main

import (
	"log"
	"net/http"
	"runtime"

	r "github.com/dancannon/gorethink"
)

// Message is a message
//type Message struct {
//	Name string      `json:"name"`
//	Data interface{} `json:"data"`
//}

// Channel is a channel
type Channel struct {
	ID   string `json:"id" gorethink:"id,omitempty"`
	Name string `json:"name" gorethink:"name"`
}

// User is a user
type User struct {
	ID   string `json:"id" gorethink:"id,omitempty"`
	Name string `json:"name" gorethink:"name"`
}

func assert(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	session, err := r.Connect(r.ConnectOpts{
		Address:  "localhost:28015",
		Database: "malice",
	})

	if err != nil {
		log.Panic(err.Error())
	}

	router := NewRouter(session)

	router.Handle("channel add", addChannel)
	router.Handle("channel subscribe", subscribeChannel)

	http.Handle("/", router)
	http.ListenAndServe(":4000", nil)
}
