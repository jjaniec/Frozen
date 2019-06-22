package main

import (
	"net"
	"fmt"
	// "strconv"
	"bufio"
	"strings"
	// "rand"
	"os"
)

type user struct {
	username string
	nickname string
	password string
	client	*connection
}

func (u *user) send_message(msg string) {
	fmt.Println("User: ", u.nickname, " sent: ", msg)
}

func (u *user) update_nickname(conn_words []string) (response string) {
	fmt.Println(conn_words)
	if len(conn_words) == 1 {
		return "No nickname specified, usage: NICK new_nickname"
	}
	new_nickname := conn_words[1]
	if len(new_nickname) == 0 || len(strings.Fields(new_nickname)) == 0 {
		return "Invalid arguments"
	}
	// Check if nickname already taken
	for _, e := range current_users {
		if (e.nickname == new_nickname) {
			return "Nickname already in use."
		}
	}

	fmt.Println("User nickname: ", u.nickname, " updated as ", new_nickname)
	u.nickname = new_nickname
	current_users = append(current_users, u)
	return
}

type server struct {
	prefix		string
	listener	net.Listener
	port		string
}

func (s *server) start() {
	srv, err := net.Listen("tcp4", s.port)
	if err != nil {
		fmt.Println(err)
		return
	}
	s.listener = srv
}

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

type channel struct {
	name	string
	subscribed_users []*user
}

var current_users = []*user{}
var current_connections = []*connection{}
var current_servers = []*server{}
var current_channels = []*channel{}

func main() {
	root_user := user{username: "root", nickname: "root", password: "toor"}
	current_users = append(current_users, &root_user)
	// users := []user{ {username: "root", nickname: "root", password: "toor"} }
	servers := []server{ {prefix: "127.0.0.1", port: ":4242"} }
	for _, e := range servers {
		e.start()
		current_servers = append(current_servers, &e)

		defer e.listener.Close() // At function end, stop server
		for {
			c, err := e.listener.Accept()
			if err != nil {
				fmt.Println(err)
				return
			}

			conn := connection{ conn: c, addr: c.RemoteAddr().String() }
			conn.session = &user{}
			current_connections = append(current_connections, &conn)
			go conn.handler()
		}
	}
}
