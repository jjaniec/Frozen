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
	server	*server
}

func (c *connection) end(reason string) {
	fmt.Println("Connection: ", c.addr, " ended w/ reason: ", reason)
	//remove connection & user from slices
	// tmp[&c] := current_connections[len(current_connections) - 1]
	// current_connections = 
	// current_users = 
}

func (c *connection) send(text string) {
	if len(text) > 0 {
		resp := fmt.Sprintf("%s\n", text)
		c.conn.Write([]byte(resp))
		fmt.Println("Answer", resp)
	}
}

func (c *connection) format_resp(args ...string) (string) {
	ret := []string {fmt.Sprintf(":%s", c.server.prefix)}
	for _, arg := range args {
		ret = append(ret, arg)
	}
	return strings.Join(ret[:], " ")
}

func (c *connection) receive() (text []string) {
	reader := bufio.NewReader(c.conn)
	var lines []string
	lines = nil
	for {
		netData, err := reader.ReadString('\n')
		lines = append(lines, netData)
		fmt.Println(lines)
		if err != nil {
			fmt.Println("err", err)
			return nil
		}
		if len(lines) > 2 {
			break
		}
	}
	fmt.Println(lines)
	return lines
}

func (c *connection) handle_line(words []string) {
	switch words[0] {
	case "ping":
		c.send("pong")
	case "kill":
		os.Exit(1)
	case "PASS":
		if len(words) < 1 {
			c.send("ERR_NEEDMOREPARAMS")
		} else {
			c.handle_cmd_pass(words[1])
		}
	case "NICK":
		if len(words) < 1 {
			c.send("ERR_NEEDMOREPARAMS")
		} else {
			resp_code, resp_str := c.handle_cmd_nick(words[1])
			if resp_str != "" && resp_code != ERR_NICKNAMEINUSE {
				c.send(c.format_resp(resp_code, resp_str))
				c.session.client = c
			} else {
				c.send(c.format_resp(resp_code, "*", words[1], resp_str))
			}
		}
	case "USER":
		if len(words) < 5 {
			c.send("ERR_NEEDMOREPARAMS") // send need more params
		} else {
			resp_code, resp_str := c.handle_cmd_user(words[1], words[2], words[3], words[4])
			c.send(c.format_resp(resp_code, c.session.username, resp_str))
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
	case "PRIVMSG":
		if len(words) < 2 {
			c.send("ERR_NEEDMOREPARAMS")
		} else {
			c.handle_cmd_privmsg(words[1], words[2])
		}
	}
}

func (c *connection) handler() {
	fmt.Printf("Serving %s\n", c.conn.RemoteAddr().String())
	defer c.end("Client lost")
	for {
		lines := c.receive()
		if lines == nil {
			print("nil")
			break
		}
		for _, line := range lines {
			temp := strings.TrimSpace(string(line))
			words := strings.Fields(temp)
			if len(words) == 0 {
				break
			}
			fmt.Println("Received: ", words, " from: ", c.addr)
			c.handle_line(words)
		}
	}
}
