package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
	"net/http"
)

type Message struct {
	Name string      `json:"name"`
	Data interface{} `json:"data"`
}

type Channel struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":4000", nil)
}
func handler(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "Hello from go")

	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		var inMesssage Message
		var outMessage Message
		if err := socket.ReadJSON(&inMesssage); err != nil {
			fmt.Println(err)
			break
		}

		fmt.Printf("%#v\n", inMesssage)
		switch inMesssage.Name {
		case "channel add":
			err := addChannel(inMesssage.Data)
			if err != nil {
				outMessage = Message{"error", err}
				if err := socket.WriteJSON(outMessage); err != nil {
					fmt.Println(err)
					break
				}
			}
		}

	}
}

func addChannel(data interface{}) error {
	var channel Channel

	err := mapstructure.Decode(data, &channel)
	if err != nil {
		return err
	}
	channel.Id = "1"
	fmt.Printf("%#v\n", channel)
	return nil
}
