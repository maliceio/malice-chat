package main

import (
	"log"
	"net/http"
	"time"

	r "github.com/dancannon/gorethink"
)

// Version stores the plugin's version
var Version string

// BuildTime stores the plugin's build time
var BuildTime string

func main() {
	session, err := r.Connect(r.ConnectOpts{
		Address:       "retrhinkdb:28015",
		Timeout:       1 * time.Second,
		MaxIdle:       3,
		MaxOpen:       10,
		DiscoverHosts: true,
		Database:      "malice",
	})

	if err != nil {
		log.Panic(err.Error())
	}

	router := NewRouter(session)

	router.Handle("channel add", addChannel)
	router.Handle("channel subscribe", subscribeChannel)
	router.Handle("channel unsubscribe", unsubscribeChannel)

	router.Handle("user edit", editUser)
	router.Handle("user subscribe", subscribeUser)
	router.Handle("user unsubscribe", unsubscribeUser)

	router.Handle("message add", addChannelMessage)
	router.Handle("message subscribe", subscribeChannelMessage)
	router.Handle("message unsubscribe", unsubscribeChannelMessage)
	http.Handle("/", router)
	http.ListenAndServe(":4000", nil)
}
