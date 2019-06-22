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

func handleConnection(c net.Conn) {
	fmt.Printf("Serving %s\n", c.RemoteAddr().String())
	defer c.Close()
	for {
			netData, err := bufio.NewReader(c).ReadString('\n')
			if err != nil {
					fmt.Println(err)
					return
			}

			temp := strings.TrimSpace(string(netData))
			if temp == "STOP" {
					break
			}

			// result := strconv.Itoa(random()) + "\n"
			c.Write([]byte("Some response"))
	}
	// c.Close()
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

		go handleConnection(c)
	}
}
