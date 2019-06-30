package main

import (
	"fmt"
	// "strconv"
)

var current_users = []*user{}
var current_connections = []*connection{}
var current_servers = []*server{}
var current_channels = []*channel{}

func main() {
	root_user := user{username: "root", nickname: "root", password: "toor"}
	current_users = append(current_users, &root_user)
	home_channel := channel{name: "#home", topic: "Base channel, owned by root", subscribed_users: []*user{&root_user}}
	current_channels = append(current_channels, &home_channel)
	servers := []server{{prefix: "127.0.0.1", port: ":4242"}}
	for _, e := range servers {
		e.start()
		current_servers = append(current_servers, &e)

		defer e.listener.Close() // At function end, stop server
		for {
			c, err := e.listener.Accept()
			if err != nil {
				fmt.Println("err", err)
				return
			}

			conn := connection{conn: c, addr: c.RemoteAddr().String()}
			conn.session = &user{}
			conn.server = &e
			current_connections = append(current_connections, &conn)
			home_channel.subscribed_users = append(home_channel.subscribed_users, conn.session)
			go conn.handler()
		}
	}
}
