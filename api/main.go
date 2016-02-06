package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	r "github.com/dancannon/gorethink"
)

// Version stores the plugin's version
var Version string

// BuildTime stores the plugin's build time
var BuildTime string

func setUpRethinkDB(session *r.Session) error {
	// Create malice DB
	resp, err := r.DBCreate("malice").RunWrite(session)
	if err != nil {
		fmt.Print(err)
	} else {
		fmt.Printf("%d DB created\n", resp.DBsCreated)
	}
	// Create channel Table
	resp, err = r.DB("malice").TableCreate("channel").RunWrite(session)
	if err != nil {
		fmt.Print(err)
	} else {
		fmt.Printf("%d Table created\n", resp.TablesCreated)
	}
	// Create message Table
	resp, err = r.DB("malice").TableCreate("message").RunWrite(session)
	if err != nil {
		fmt.Print(err)
	} else {
		resp, err = r.DB("malice").Table("message").IndexCreate("createdAt", r.IndexCreateOpts{
			Multi: true,
		}).RunWrite(session)
		if err != nil {
			fmt.Print(err)
		}
		resp, err = r.DB("malice").Table("message").IndexCreate("channelId", r.IndexCreateOpts{
			Multi: true,
		}).RunWrite(session)
		if err != nil {
			fmt.Print(err)
		}
		fmt.Printf("%d Table created\n", resp.TablesCreated)
	}
	// Create user Table
	resp, err = r.DB("malice").TableCreate("user").RunWrite(session)
	if err != nil {
		fmt.Print(err)
	} else {
		fmt.Printf("%d Table created\n", resp.TablesCreated)
	}

	session.Use("malice")

	return err
}

func main() {
	// addrs, err := net.LookupHost("rethinkdb")
	// if err != nil {
	// 	log.Panic(err.Error())
	// }
	// rethinkAddr := addrs[0] + ":28015"
	session, err := r.Connect(r.ConnectOpts{
		// Address: "localhost:28015",
		// Address: "192.168.99.100:28015",
		Address: "db:28015",
		// Address: rethinkAddr,
		Timeout:  5 * time.Second,
		Database: "malice",
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
