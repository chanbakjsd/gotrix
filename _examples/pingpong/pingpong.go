package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/chanbakjsd/gotrix"
	"github.com/chanbakjsd/gotrix/event"
)

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func handleMessage(c *gotrix.Client, m event.RoomMessage) {
	// If it's a notice (another bot's message) or not "ping", ignore.
	if m.MsgType == event.RoomMessageNotice || m.Body != "ping" {
		return
	}
	// Otherwise, send pong!
	_, err := c.SendNotice(m.RoomID, "Pong!")
	panicIfErr(err)
}

var username = flag.String("user", "", "username")
var password = flag.String("pass", "", "password")
var url = flag.String("url", "http://localhost:8008", "url")

func main() {
	flag.Parse()
	// Ask for username and password.
	if *username == "" || *password == "" {
		fmt.Printf("user and pass flags not provided.")
		return
	}

	// Construct the client.
	cli, err := gotrix.New(*url)
	panicIfErr(err)

	// Login using provided creds.
	panicIfErr(cli.LoginPassword(*username, *password))

	// Register the handler.
	cli.AddHandler(handleMessage)

	// Start the connection.
	panicIfErr(cli.Open())

	// Wait until interrupt happens.
	fmt.Println("Ctrl-C to terminate program.")
	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, os.Interrupt)
	<-terminate

	// Close the connection and logout as we're not persisting the session.
	panicIfErr(cli.Close())
	panicIfErr(cli.Logout())
}
