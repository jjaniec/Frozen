package main

import (
	"net"
	"bufio"
	"strings"
	"fmt"
	"os"
)

type connection struct {
	conn	net.Conn
	addr	string
	session	*user
}

func (c *connection) end(reason string) {
	fmt.Println("Connection: ", c.addr, " ended w/ reason: ", reason)
	//remove connection & user from slices
	// tmp[&c] := current_connections[len(current_connections) - 1]
	// current_connections = 
	// current_users = 
}

func (c *connection) receive() (text []string) {
	netData, err := bufio.NewReader(c.conn).ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return nil
	}

	temp := strings.TrimSpace(string(netData))
	words := strings.Fields(temp)
	return words
}

func (c *connection) send(text string) {
	if len(text) > 0 {
		resp := fmt.Sprintf("%s\n", text)
		c.conn.Write([]byte(resp))
	}
}

func (c *connection) handler() {
	fmt.Printf("Serving %s\n", c.conn.RemoteAddr().String())
	defer c.end("Client lost")
	for {
		words := c.receive()
		if len(words) == 0 {
			break
		}
		fmt.Println("Received: ", words, " from: ", c.addr)
		switch words[0] {
		case "ping":
			c.send("pong")
		case "kill":
			os.Exit(1)
		case "NICK":
			resp := c.session.update_nickname(words)
			c.send(resp)
			if len(resp) > 0 {
				c.session.client = c
			}
		case "NAMES":
			fmt.Println(current_connections)
			users_connected_count := 0
			for _, e := range current_connections {
				if (e.session != nil) {
					c.send(e.session.username)
					users_connected_count++
				}
			}
			if (users_connected_count == 0) {
				c.send("No users currently connected")
			}
		case "JOIN":
			if (*c.session == user{}) {
				c.send("You must login first !")
			}
		}
	}
}
