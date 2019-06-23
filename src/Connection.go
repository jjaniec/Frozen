package main

import (
	"net"
	"strings"
	"fmt"
	"os"
	"io"
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
			if resp_code == ERR_NICKNAMEINUSE {
				c.send(c.format_resp(resp_code, "*", words[1], resp_str))
			} else if resp_str != "" && resp_code != ERR_NICKNAMEINUSE {
				c.send(c.format_resp(resp_code, resp_str))
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
		}
	case "PRIVMSG":
		if len(words) < 2 {
			c.send("ERR_NEEDMOREPARAMS")
		} else {
			c.handle_cmd_privmsg(words[1], words[2])
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
	fmt.Println("total size:", len(buf), " buf: ", buf)
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
			c.handle_line(words)
		}
	}
}
