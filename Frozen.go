package main

import (
	"net"
	"fmt"
	// "strconv"
	"bufio"
	"strings"
	// "rand"
)

type user struct {
	username string
	nickname string
	password string
}

func (u *user) send_message(msg string) {
	fmt.Println("User: ", u.nickname, " sent: ", msg)
}

type server struct {
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
}

func (c *connection) end(reason string) {
	fmt.Println("Connection: ", c.addr, " ended w/ reason: ", reason)
}

func (c *connection) receive() (text *string) {
	netData, err := bufio.NewReader(c.conn).ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return nil
	}

	temp := strings.TrimSpace(string(netData))
	if temp == "STOP" {
		c.end("Client goodbye received")
	}
	return &temp
}

func (c *connection) handler() {
	fmt.Printf("Serving %s\n", c.conn.RemoteAddr().String())
	defer c.end("Client lost")
	for {
		text := c.receive()
		if text == nil {
			break
		}
		fmt.Println("Received: ", *text, " from: ", c.addr)
		c.conn.Write([]byte("Some response"))
	}
}

func main() {
	// channels := make(chan string)
	users := []user{ {username: "root", nickname: "root", password: "toor"} }
	users[0].send_message("Test")

	srv := server{port: ":4242"}
	srv.start()
	defer srv.listener.Close() // At function end, stop server
	for {
		c, err := srv.listener.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}

		conn := connection{ conn: c, addr: c.RemoteAddr().String() }
		go conn.handler()
		// go handleConnection({})
	}
}
