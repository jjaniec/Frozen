package main

import (
	"net"
	"strings"
	"fmt"
	// "os"
	"io"
)

const ERR_NORECIPIENT = "411"
const ERR_NOTEXTTOSEND = "412"
const ERR_NOTREGISTERED = "451" // ":You have not registered"

type connection struct {
	conn	net.Conn
	addr	string
	session	*user
	server	*server
}

func (c *connection) end(reason string) {
	fmt.Println("Connection: ", c.addr, " ended w/ reason: ", reason)
	if (c.session != nil && c.session.nickname != "") {
		for i, e := range current_users {
			if (e.nickname == c.session.nickname) {
				fmt.Println("Delete user", c.session.nickname)
				current_users[i] = current_users[len(current_users) - 1]
				current_users[len(current_users) - 1] = nil
				current_users = current_users[:len(current_users) - 1]
			}
		}
	}
	for i, e := range current_connections {
		if (e.addr == c.addr) {
			fmt.Println("Delete connection", c)
			current_connections[i] = current_connections[len(current_connections) - 1]
			current_connections[len(current_connections) - 1] = nil
			current_connections = current_connections[:len(current_connections) - 1]

		}
	}
}

func (c *connection) send(text string) {
	if len(text) > 0 {
		resp := fmt.Sprintf("%s\r\n", text)
		c.conn.Write([]byte(resp))
		fmt.Println("Answer", resp)
	}
}

func (c *connection) format_resp(args ...string) (string) {
	var ret []string
	if (args[0][0] != ':') {
		ret = []string {fmt.Sprintf(":%s", c.server.prefix)}
	}
	// ret := []string {fmt.Sprintf(":%s", c.server.prefix)}
	for _, arg := range args {
		ret = append(ret, arg)
	}
	return strings.Join(ret[:], " ")
}

func (c *connection) handle_line(words []string, raw_line string) {
	// Here,handle nickname suffixes (like :Bob PRIVMSG Alex :bla bla)
	cmd_str := words[0]
	switch cmd_str {
	case "PING":
		c.send(c.format_resp("PONG", "ping.frozen", fmt.Sprintf(":%s", c.server.prefix)))
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
			if resp_code == ERR_NICKNAMEINUSE {
				c.send(c.format_resp(resp_code, "*", words[1], resp_str))
			} else if resp_code != ERR_NICKNAMEINUSE {
				// c.send(c.format_resp(resp_code, resp_str))
				c.session.client = c
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
		var nicknames_fmt []string

		if (c.session.nickname == "") {
			return
		}
		if (len(words) > 1) {
			nicknames_fmt = c.handle_cmd_names(words[1])
		} else {
			nicknames_fmt = c.handle_cmd_names("")
		}
		for _, e := range nicknames_fmt {
			c.send(c.format_resp(RPL_NAMREPLY, c.session.nickname, e))
		}
		c.send(c.format_resp(RPL_ENDOFNAMES, c.session.nickname, "*", ":End of /NAMES list."))
	case "JOIN":
		if (*c.session == user{}) {
			c.send("You must login first !")
		} else {
			if len(words) < 2 {
				c.send(c.format_resp("461", c.session.nickname, "JOIN", ":Not enough parameters"))
			} else {
				c.handle_cmd_join(words[1])
			}
		}
	case "LIST":
		c.handle_cmd_list()
	case "TOPIC":
		if len(words) < 2 {
			c.send("ERR_NEEDMOREPARAMS") // send need more params
		} else {
			if (len(words) == 2) {
				c.handle_cmd_topic(words[1], "")
			} else {
				c.handle_cmd_topic(words[1], strings.Join(words[2:], " "))
			}
		}
	case "PART":
		if len(words) < 1 {
			c.send("ERR_NEEDMOREPARAMS") // send need more params
		} else {
			c.handle_cmd_part(words[1])
		}
	case "PRIVMSG":
		// print(words)
		fmt.Println(c.session.client)
		if (c.session.nickname == "") {
			c.send(c.format_resp(ERR_NOTREGISTERED, "*", ":You have not registered"))
		} else if len(words) == 1 {
			c.send(c.format_resp(ERR_NORECIPIENT, c.session.nickname, ":No recipient given (PRIVMSG)"))
		} else if len(words) == 2 {
			c.send(c.format_resp(ERR_NOTEXTTOSEND, c.session.nickname, ":No text to send"))
		} else {
			c.handle_cmd_privmsg(words[1], raw_line)
		}
	}
}

func (c *connection) receive() (status error, text []string) {
	// https://stackoverflow.com/questions/24339660/read-whole-data-with-golang-net-conn-read
	buf := make([]byte, 0, 4096)
	tmp := make([]byte, 1024)
	n, err := c.conn.Read(tmp)
	if err != nil {
		if err != io.EOF {
			fmt.Println("read error:", err)
		}
		return err, nil
	}
	buf = append(buf, tmp[:n]...)
	// fmt.Println("Append: ", tmp)
	// fmt.Println("total size:", len(buf), " buf: ", buf)
	lines := strings.Split(string(buf), "\r\n")
	return nil, lines
}

func (c *connection) handler() {
	fmt.Printf("Serving %s\n", c.conn.RemoteAddr().String())
	defer c.end("Client lost")
	for {
		status, lines := c.receive()
		if status != nil {
			break
		}
		print(lines)
		for _, line := range lines {
			temp := strings.TrimSpace(string(line))
			words := strings.Fields(temp)
			if len(words) == 0 {
				break
			}
			fmt.Println("Received: ", words, " from: ", c.addr)
			c.handle_line(words, line)
		}
	}
}
