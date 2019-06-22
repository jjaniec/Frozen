package main

import (
	"fmt"
	"strings"
)

const RPL_WELCOME = "001"
const ERR_NICKNAMEINUSE = "433"
const ERR_NONICKNAMEGIVEN = "431"
const ERR_NOSUCHNICK = "401"

func (c *connection) handle_cmd_pass(password string) {
	// https://tools.ietf.org/html/rfc1459#section-4.1.1
	c.session.password = password
}

func (c *connection) handle_cmd_nick(nickname string) (resp_code string, resp_str string) {
	// https://tools.ietf.org/html/rfc1459#section-4.1.2
	// Check if nickname already taken
	for _, e := range current_users {
		if (e.nickname == nickname) {
			return ERR_NICKNAMEINUSE, "Nickname is already in use."
		}
	}
	fmt.Println("User nickname: ", c.session.nickname, " updated as ", nickname)
	c.session.nickname = nickname
	current_users = append(current_users, c.session)
	return
}

func (c *connection) handle_cmd_user(username string, hostname string, servername string, realname string) (resp_code string, resp_str string){
	// https://tools.ietf.org/html/rfc1459#section-4.1.3
	// Only nicknames must be unique
	if (c.session.nickname == "") {
		return ERR_NONICKNAMEGIVEN, fmt.Sprintf(":No nickname given")
	}
	c.session.username = username
	c.session.realname = realname
	return RPL_WELCOME, fmt.Sprintf(":Welcome to the Internet Relay Network %s!%s@%s", c.session.nickname, c.session.username, c.server.prefix)
}

func (c *connection) handle_cmd_privmsg(receiver string, text string) {
	// https://tools.ietf.org/html/rfc1459#section-4.4.1
	receivers := strings.Split(receiver, ",")
	for _, e := range receivers {
		if (e[0] != '#' && e[0] != '$') {
			for _, u := range current_users {
				if (u.nickname == e) {
					u.client.send(fmt.Sprintf(":%s!%s@%s PRIVMSG %s :%s", c.session.nickname, c.session.username, c.server.prefix, receiver, text))
					return
				}
			}
			c.send(c.format_resp(ERR_NOSUCHNICK, ":No such nick/channel"))
		} else {
			// Handle host / server mask
			print("TODO")
		}
	}
}
