package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	r "github.com/dancannon/gorethink"
)

// Version stores the plugin's version
var Version string

// BuildTime stores the plugin's build time
var BuildTime string

func setUpRethinkDB(session *r.Session) error {
	resp, err := r.DBCreate("malice").RunWrite(session)
	if err != nil {
		fmt.Print(err)
	}
	fmt.Printf("%d DB created", resp.DBsCreated)
	r.TableCreate("channel", r.TableCreateOpts{PrimaryKey: "channelId"})
	fmt.Printf("%d Table created", resp.DBsCreated)
	r.TableCreate("message")
	fmt.Printf("%d Table created", resp.DBsCreated)
	r.TableCreate("user")
	fmt.Printf("%d Table created", resp.DBsCreated)

	session.Use("malice")

	return err
}

func main() {
	addrs, err := net.LookupHost("rethinkdb")
	if err != nil {
		log.Panic(err.Error())
	}
	rethinkAddr := addrs[0] + ":28015"
	session, err := r.Connect(r.ConnectOpts{
		Address: rethinkAddr,
		Timeout: 5 * time.Second,
	})

	if err != nil {
		log.Panic(err.Error())
	}

	err = setUpRethinkDB(session)
	if err != nil {
		fmt.Print(err)
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
